package service

import (
	"context"
	"fmt"
	"time"

	"starledger/ent"
	"starledger/ent/bill"
)

type BillService struct {
	client *ent.Client
}

func NewBillService(client *ent.Client) *BillService {
	return &BillService{client: client}
}

func (s *BillService) List(ctx context.Context, tenantID, page, pageSize int, billType, status string) ([]*ent.Bill, int, error) {
	offset := (page - 1) * pageSize
	q := s.client.Bill.Query().Where(bill.TenantID(tenantID))

	if billType != "" {
		q = q.Where(bill.BillType(billType))
	}
	if status != "" {
		q = q.Where(bill.PaymentStatusEQ(bill.PaymentStatus(status)))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := q.
		Offset(offset).
		Limit(pageSize).
		Order(ent.Desc(bill.FieldCreatedAt)).
		All(ctx)
	return items, total, err
}

func (s *BillService) Get(ctx context.Context, id, tenantID int) (*ent.Bill, error) {
	return s.client.Bill.Query().
		Where(bill.ID(id), bill.TenantID(tenantID)).
		Only(ctx)
}

func (s *BillService) Create(ctx context.Context, tenantID int, billType string, resourceID *int, amount float64, billDate, dueDate time.Time, remark string) (*ent.Bill, error) {
	billNo := fmt.Sprintf("BL%s%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)

	creator := s.client.Bill.Create().
		SetTenantID(tenantID).
		SetBillNo(billNo).
		SetBillType(billType).
		SetAmount(amount).
		SetBillDate(billDate).
		SetDueDate(dueDate).
		SetPaymentStatus(bill.PaymentStatusPending).
		SetRemark(remark)

	if resourceID != nil {
		creator = creator.SetRelatedResourceID(*resourceID)
	}

	return creator.Save(ctx)
}

func (s *BillService) Pay(ctx context.Context, id, tenantID int, paymentMethod, invoiceNo string) (*ent.Bill, error) {
	return s.client.Bill.UpdateOneID(id).
		SetPaymentStatus(bill.PaymentStatusPaid).
		SetPaidDate(time.Now()).
		SetPaymentMethod(paymentMethod).
		SetInvoiceNo(invoiceNo).
		Save(ctx)
}

func (s *BillService) Cancel(ctx context.Context, id, tenantID int) (*ent.Bill, error) {
	return s.client.Bill.UpdateOneID(id).
		SetPaymentStatus(bill.PaymentStatusCancelled).
		Save(ctx)
}

func (s *BillService) Overdue(ctx context.Context, tenantID int) ([]*ent.Bill, error) {
	return s.client.Bill.Query().
		Where(
			bill.TenantID(tenantID),
			bill.PaymentStatusEQ(bill.PaymentStatusOverdue),
		).
		Order(ent.Asc(bill.FieldDueDate)).
		All(ctx)
}

func (s *BillService) Summary(ctx context.Context, tenantID int) (map[string]interface{}, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	totalAmount, _ := s.client.Bill.Query().
		Where(bill.TenantID(tenantID), bill.BillDateGTE(startOfMonth)).
		Aggregate(ent.Sum(bill.FieldAmount)).
		Int(ctx)

	pendingCount, _ := s.client.Bill.Query().
		Where(bill.TenantID(tenantID), bill.PaymentStatusEQ(bill.PaymentStatusPending)).
		Count(ctx)

	paidCount, _ := s.client.Bill.Query().
		Where(bill.TenantID(tenantID), bill.PaymentStatusEQ(bill.PaymentStatusPaid)).
		Count(ctx)

	overdueCount, _ := s.client.Bill.Query().
		Where(bill.TenantID(tenantID), bill.PaymentStatusEQ(bill.PaymentStatusOverdue)).
		Count(ctx)

	return map[string]interface{}{
		"monthly_total":  totalAmount,
		"pending_count":  pendingCount,
		"paid_count":     paidCount,
		"overdue_count":  overdueCount,
	}, nil
}

// CheckOverdue updates overdue bills status.
func (s *BillService) CheckOverdue(ctx context.Context) error {
	now := time.Now()
	_, err := s.client.Bill.Update().
		Where(
			bill.PaymentStatusEQ(bill.PaymentStatusPending),
			bill.DueDateLTE(now),
		).
		SetPaymentStatus(bill.PaymentStatusOverdue).
		Save(ctx)
	return err
}
