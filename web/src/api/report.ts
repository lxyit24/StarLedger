import request from '../utils/request'

export function reportOverview() {
  return request.get('/reports/overview')
}

export function monthlyTrend() {
  return request.get('/reports/monthly-trend')
}

export function billTypeDistribution() {
  return request.get('/reports/bill-type')
}

export function billStatusSummary() {
  return request.get('/reports/bill-status')
}

export function serverCostAnalysis() {
  return request.get('/reports/server-cost')
}

export function taskStatistics() {
  return request.get('/reports/task-stats')
}
