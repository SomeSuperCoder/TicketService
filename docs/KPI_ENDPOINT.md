# KPI Monitoring Endpoint

## Endpoint
`GET /api/v1/monitoring/kpi`

## Query Parameters
- `period` (optional): Time period for metrics calculation
  - Values: `week`, `month`, `year`
  - Default: `month`

## Response
```json
{
  "avg_response_days": 2.5,
  "overdue_count": 15,
  "satisfaction_index": 75.5
}
```

## Metrics Explained

### Average Response Days
- Calculates the average time (in days) from ticket creation to first comment
- Only includes tickets created within the specified period
- Only counts tickets that have received at least one response

### Overdue Count
- Total number of currently overdue tickets (across all time, not just the period)
- Uses the centralized `v_ticket_overdue_status` view
- Tickets are overdue if they've been in open/init status for more than 7 days

### Satisfaction Index
- A calculated metric based on ticket resolution performance
- Formula: `(closed_tickets + on_time_tickets * 0.5) / total_tickets * 100`
- Range: 0-100
- Higher values indicate better performance
- Considers:
  - Closed tickets (100% weight)
  - On-time open tickets (50% weight)
  - Only tickets within the specified period

## Period Mapping
- `week`: Last 7 days
- `month`: Last 30 days
- `year`: Last 365 days

## Example Usage
```bash
# Get KPI for the last week
GET /api/v1/monitoring/kpi?period=week

# Get KPI for the last month (default)
GET /api/v1/monitoring/kpi

# Get KPI for the last year
GET /api/v1/monitoring/kpi?period=year
```
