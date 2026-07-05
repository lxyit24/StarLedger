package service

import (
	"context"
	"fmt"
	"time"

	"starledger/ent"
	"starledger/ent/invoice"
)

type InvoiceService struct {
	client *ent.Client
}

func NewInvoiceService(client *ent.Client) *InvoiceService {
	return &InvoiceService{client: client}
}

// List returns paginated invoices with optional filters.
func (s *InvoiceService) List(ctx context.Context, tenantID, page, pageSize int, status, invoiceType string) ([]*ent.Invoice, int, error) {
	offset := (page - 1) * pageSize
	q := s.client.Invoice.Query().Where(invoice.TenantID(tenantID))

	if status != "" {
		q = q.Where(invoice.StatusEQ(invoice.Status(status)))
	}
	if invoiceType != "" {
		q = q.Where(invoice.InvoiceTypeEQ(invoice.InvoiceType(invoiceType)))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := q.
		Offset(offset).
		Limit(pageSize).
		Order(ent.Desc(invoice.FieldCreatedAt)).
		All(ctx)
	return items, total, err
}

// Get returns a single invoice by ID.
func (s *InvoiceService) Get(ctx context.Context, id, tenantID int) (*ent.Invoice, error) {
	return s.client.Invoice.Query().
		Where(invoice.ID(id), invoice.TenantID(tenantID)).
		WithBill().
		Only(ctx)
}

// Create creates a new invoice in draft status.
func (s *InvoiceService) Create(ctx context.Context, tenantID int, req CreateInvoiceReq) (*ent.Invoice, error) {
	invoiceNo := req.InvoiceNo
	if invoiceNo == "" {
		invoiceNo = fmt.Sprintf("INV%s%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
	}

	totalAmount := req.Amount + req.TaxAmount

	creator := s.client.Invoice.Create().
		SetTenantID(tenantID).
		SetInvoiceNo(invoiceNo).
		SetInvoiceCode(req.InvoiceCode).
		SetInvoiceType(invoice.InvoiceType(req.InvoiceType)).
		SetAmount(req.Amount).
		SetTaxAmount(req.TaxAmount).
		SetTotalAmount(totalAmount).
		SetTaxRate(req.TaxRate).
		SetBuyerName(req.BuyerName).
		SetBuyerTaxNo(req.BuyerTaxNo).
		SetBuyerAddress(req.BuyerAddress).
		SetBuyerBank(req.BuyerBank).
		SetSellerName(req.SellerName).
		SetSellerTaxNo(req.SellerTaxNo).
		SetSellerAddress(req.SellerAddress).
		SetSellerBank(req.SellerBank).
		SetItemsDetail(req.ItemsDetail).
		SetRemark(req.Remark).
		SetStatus(invoice.StatusDraft)

	if req.BillID > 0 {
		creator = creator.SetBillID(req.BillID)
	}
	if !req.InvoiceDate.IsZero() {
		creator = creator.SetInvoiceDate(req.InvoiceDate)
	}

	return creator.Save(ctx)
}

// Update updates a draft invoice.
func (s *InvoiceService) Update(ctx context.Context, id, tenantID int, req UpdateInvoiceReq) (*ent.Invoice, error) {
	inv, err := s.client.Invoice.Query().
		Where(invoice.ID(id), invoice.TenantID(tenantID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	if inv.Status != invoice.StatusDraft {
		return nil, fmt.Errorf("只能编辑草稿状态的发票")
	}

	updater := s.client.Invoice.UpdateOneID(id).
		SetInvoiceCode(req.InvoiceCode).
		SetAmount(req.Amount).
		SetTaxAmount(req.TaxAmount).
		SetTotalAmount(req.Amount + req.TaxAmount).
		SetTaxRate(req.TaxRate).
		SetBuyerName(req.BuyerName).
		SetBuyerTaxNo(req.BuyerTaxNo).
		SetBuyerAddress(req.BuyerAddress).
		SetBuyerBank(req.BuyerBank).
		SetSellerName(req.SellerName).
		SetSellerTaxNo(req.SellerTaxNo).
		SetSellerAddress(req.SellerAddress).
		SetSellerBank(req.SellerBank).
		SetItemsDetail(req.ItemsDetail).
		SetRemark(req.Remark)

	return updater.Save(ctx)
}

// Issue issues a draft invoice (changes status to issued).
func (s *InvoiceService) Issue(ctx context.Context, id, tenantID int) (*ent.Invoice, error) {
	inv, err := s.client.Invoice.Query().
		Where(invoice.ID(id), invoice.TenantID(tenantID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	if inv.Status != invoice.StatusDraft {
		return nil, fmt.Errorf("只能开具草稿状态的发票")
	}

	return s.client.Invoice.UpdateOneID(id).
		SetStatus(invoice.StatusIssued).
		SetInvoiceDate(time.Now()).
		Save(ctx)
}

// Cancel cancels an issued invoice.
func (s *InvoiceService) Cancel(ctx context.Context, id, tenantID int) (*ent.Invoice, error) {
	inv, err := s.client.Invoice.Query().
		Where(invoice.ID(id), invoice.TenantID(tenantID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	if inv.Status != invoice.StatusIssued {
		return nil, fmt.Errorf("只能作废已开具的发票")
	}

	return s.client.Invoice.UpdateOneID(id).
		SetStatus(invoice.StatusCancelled).
		Save(ctx)
}

// Delete deletes a draft invoice.
func (s *InvoiceService) Delete(ctx context.Context, id, tenantID int) error {
	inv, err := s.client.Invoice.Query().
		Where(invoice.ID(id), invoice.TenantID(tenantID)).
		Only(ctx)
	if err != nil {
		return err
	}
	if inv.Status != invoice.StatusDraft {
		return fmt.Errorf("只能删除草稿状态的发票")
	}

	return s.client.Invoice.DeleteOneID(id).Exec(ctx)
}

// Summary returns invoice statistics.
func (s *InvoiceService) Summary(ctx context.Context, tenantID int) (map[string]interface{}, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// This month issued invoices
	issuedCount, _ := s.client.Invoice.Query().
		Where(
			invoice.TenantID(tenantID),
			invoice.StatusEQ(invoice.StatusIssued),
			invoice.InvoiceDateGTE(startOfMonth),
		).Count(ctx)

	// This month total amount
	issuedInvoices, _ := s.client.Invoice.Query().
		Where(
			invoice.TenantID(tenantID),
			invoice.StatusEQ(invoice.StatusIssued),
			invoice.InvoiceDateGTE(startOfMonth),
		).All(ctx)

	var monthlyTotal float64
	for _, inv := range issuedInvoices {
		monthlyTotal += inv.TotalAmount
	}

	// Draft count
	draftCount, _ := s.client.Invoice.Query().
		Where(invoice.TenantID(tenantID), invoice.StatusEQ(invoice.StatusDraft)).
		Count(ctx)

	// Cancelled count
	cancelledCount, _ := s.client.Invoice.Query().
		Where(invoice.TenantID(tenantID), invoice.StatusEQ(invoice.StatusCancelled)).
		Count(ctx)

	// Total count
	totalCount, _ := s.client.Invoice.Query().
		Where(invoice.TenantID(tenantID)).
		Count(ctx)

	return map[string]interface{}{
		"monthly_issued_count": issuedCount,
		"monthly_total_amount": monthlyTotal,
		"draft_count":         draftCount,
		"cancelled_count":     cancelledCount,
		"total_count":         totalCount,
	}, nil
}

// Request types

type CreateInvoiceReq struct {
	InvoiceNo    string  `json:"invoice_no"`
	InvoiceCode  string  `json:"invoice_code"`
	InvoiceType  string  `json:"invoice_type"`
	BillID       int     `json:"bill_id"`
	Amount       float64 `json:"amount"`
	TaxAmount    float64 `json:"tax_amount"`
	TaxRate      float64 `json:"tax_rate"`
	BuyerName    string  `json:"buyer_name"`
	BuyerTaxNo   string  `json:"buyer_tax_no"`
	BuyerAddress string  `json:"buyer_address"`
	BuyerBank    string  `json:"buyer_bank"`
	SellerName   string  `json:"seller_name"`
	SellerTaxNo  string  `json:"seller_tax_no"`
	SellerAddress string `json:"seller_address"`
	SellerBank   string  `json:"seller_bank"`
	InvoiceDate  time.Time `json:"invoice_date"`
	ItemsDetail  string  `json:"items_detail"`
	Remark       string  `json:"remark"`
}

type UpdateInvoiceReq struct {
	InvoiceCode  string  `json:"invoice_code"`
	Amount       float64 `json:"amount"`
	TaxAmount    float64 `json:"tax_amount"`
	TaxRate      float64 `json:"tax_rate"`
	BuyerName    string  `json:"buyer_name"`
	BuyerTaxNo   string  `json:"buyer_tax_no"`
	BuyerAddress string  `json:"buyer_address"`
	BuyerBank    string  `json:"buyer_bank"`
	SellerName   string  `json:"seller_name"`
	SellerTaxNo  string  `json:"seller_tax_no"`
	SellerAddress string `json:"seller_address"`
	SellerBank   string  `json:"seller_bank"`
	ItemsDetail  string  `json:"items_detail"`
	Remark       string  `json:"remark"`
}
