# Heatmap Endpoints

## Overview

The heatmap endpoints provide geographic visualization data for tickets based on their complaint locations. These endpoints support filtering by time period and categories.

## Endpoints

### 1. Get Heatmap Points

```
GET /api/v1/heatmap/points
```

Returns geographic points for rendering on a map with intensity values.

#### Query Parameters

- `period` (optional) - Time period: `week`, `month`, or `year` (default: `month`)
- `categories` (optional) - Array of subcategory IDs to filter by

#### Response

```json
{
  "points": [
    {
      "lat": 37.7749,
      "lng": -122.4194,
      "intensity": 15,
      "address": "POINT(-122.4194 37.7749)",
      "count": 15
    },
    {
      "lat": 37.7849,
      "lng": -122.4094,
      "intensity": 8,
      "address": "POINT(-122.4094 37.7849)",
      "count": 8
    }
  ]
}
```

#### Field Descriptions

- `lat` - Latitude coordinate
- `lng` - Longitude coordinate
- `intensity` - Heat intensity value (same as count, used for heatmap visualization)
- `address` - Geographic point in text format
- `count` - Number of tickets at this location

#### Examples

Get all points for the last month:
```
GET /api/v1/heatmap/points
```

Get points for the last week:
```
GET /api/v1/heatmap/points?period=week
```

Get points filtered by categories:
```
GET /api/v1/heatmap/points?period=month&categories=1&categories=3&categories=5
```

### 2. Get Heatmap Statistics

```
GET /api/v1/heatmap/stats
```

Returns statistics about the heatmap including top problem locations.

#### Response

```json
{
  "top_locations": [
    {
      "lat": 37.7749,
      "lng": -122.4194,
      "ticket_count": 25,
      "address": "POINT(-122.4194 37.7749)"
    },
    {
      "lat": 37.7849,
      "lng": -122.4094,
      "ticket_count": 18,
      "address": "POINT(-122.4094 37.7849)"
    }
  ],
  "total_locations": 150,
  "avg_tickets_per_location": 3.5
}
```

#### Field Descriptions

- `top_locations` - Top 5 locations with the most tickets
  - `lat` - Latitude coordinate
  - `lng` - Longitude coordinate
  - `ticket_count` - Number of tickets at this location
  - `address` - Geographic point in text format
- `total_locations` - Total number of unique locations with tickets
- `avg_tickets_per_location` - Average number of tickets per location

## Frontend Integration

### Heatmap Visualization

Use the `/heatmap/points` endpoint with a heatmap library like Leaflet.heat or Google Maps Heatmap Layer:

```javascript
// Fetch heatmap data
const response = await fetch('/api/v1/heatmap/points?period=month');
const data = await response.json();

// Convert to heatmap format
const heatmapData = data.points.map(point => ({
  lat: point.lat,
  lng: point.lng,
  intensity: point.intensity
}));

// Render with Leaflet.heat
const heat = L.heatLayer(heatmapData, {
  radius: 25,
  blur: 15,
  maxZoom: 17
}).addTo(map);
```

### Problem Locations Markers

Use the `/heatmap/stats` endpoint to show top problem locations:

```javascript
const response = await fetch('/api/v1/heatmap/stats');
const data = await response.json();

// Add markers for top locations
data.top_locations.forEach(location => {
  L.marker([location.lat, location.lng])
    .bindPopup(`${location.ticket_count} tickets`)
    .addTo(map);
});
```

## Implementation Details

- Uses PostGIS `ST_X()` and `ST_Y()` functions to extract coordinates
- Groups tickets by exact geographic location
- Only includes non-deleted and non-hidden tickets
- Points are ordered by count (descending) for performance
- Top locations limited to 5 results
- Supports filtering by subcategory IDs
- Time periods calculated from current date backwards

## Database Requirements

Requires PostGIS extension for geographic functions:
- `ST_X()` - Extract longitude
- `ST_Y()` - Extract latitude
- Geographic data stored in `complaint_details.geo_location` field
