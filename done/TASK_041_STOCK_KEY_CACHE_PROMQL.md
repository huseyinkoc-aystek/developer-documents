# TASK 041 — Stock Key Cache Metrics (PromQL Rehberi)

## Kapsam
Bu doküman, `ops_stock_key_lookup_cache_total{module,mode,result}` metriği için `/ops/metrics` tabanlı Prometheus sorgularını içerir.

Metric:
- `ops_stock_key_lookup_cache_total`
- Labels:
  - `module`: `stock_inventory|stock_sales|stock_returns|stock_subscriptions|stock_vehicle_assets`
  - `mode`: `exact|fuzzy`
  - `result`: `hit|miss`

## 1) Global Hit Ratio
Son 5 dakikada global cache hit oranı:

```promql
sum(rate(ops_stock_key_lookup_cache_total{result="hit"}[5m]))
/
clamp_min(sum(rate(ops_stock_key_lookup_cache_total[5m])), 1e-9)
```

Yüzde formatı:

```promql
100 *
sum(rate(ops_stock_key_lookup_cache_total{result="hit"}[5m]))
/
clamp_min(sum(rate(ops_stock_key_lookup_cache_total[5m])), 1e-9)
```

## 2) Module Bazlı Hit Ratio
Son 5 dakikada modül bazlı hit oranı:

```promql
100 *
sum by (module) (rate(ops_stock_key_lookup_cache_total{result="hit"}[5m]))
/
clamp_min(sum by (module) (rate(ops_stock_key_lookup_cache_total[5m])), 1e-9)
```

## 3) Module Bazlı Miss Trendi
Miss/saniye trendi (5m rate):

```promql
sum by (module) (rate(ops_stock_key_lookup_cache_total{result="miss"}[5m]))
```

Miss adedi (15m increase):

```promql
sum by (module) (increase(ops_stock_key_lookup_cache_total{result="miss"}[15m]))
```

## 4) Exact vs Fuzzy Karşılaştırma
Mode bazlı toplam lookup yoğunluğu:

```promql
sum by (mode) (rate(ops_stock_key_lookup_cache_total[5m]))
```

Exact mod hit oranı:

```promql
100 *
sum(rate(ops_stock_key_lookup_cache_total{mode="exact",result="hit"}[5m]))
/
clamp_min(sum(rate(ops_stock_key_lookup_cache_total{mode="exact"}[5m])), 1e-9)
```

Fuzzy mod hit oranı:

```promql
100 *
sum(rate(ops_stock_key_lookup_cache_total{mode="fuzzy",result="hit"}[5m]))
/
clamp_min(sum(rate(ops_stock_key_lookup_cache_total{mode="fuzzy"}[5m])), 1e-9)
```

Not: Fuzzy mod doğal olarak daha çok miss üretebilir; panelde exact/fuzzy ayrı izlenmelidir.

## 5) Dashboard Kartları (Öneri)
Grafana için önerilen 4 kart:
1. Global hit ratio (%)
2. Module hit ratio (%)
3. Miss trend by module
4. Exact/Fuzzy request share

## 6) Alert Örnekleri
### A) Global hit ratio çok düşük
Koşul (10 dk boyunca < %60):

```promql
(
  100 *
  sum(rate(ops_stock_key_lookup_cache_total{result="hit"}[10m]))
  /
  clamp_min(sum(rate(ops_stock_key_lookup_cache_total[10m])), 1e-9)
) < 60
```

### B) Bir modülde miss patlaması
Koşul (15 dk’da miss > 1000):

```promql
sum by (module) (increase(ops_stock_key_lookup_cache_total{result="miss"}[15m])) > 1000
```

### C) Exact lookup verimsizliği
Koşul (exact hit ratio < %50):

```promql
(
  100 *
  sum(rate(ops_stock_key_lookup_cache_total{mode="exact",result="hit"}[10m]))
  /
  clamp_min(sum(rate(ops_stock_key_lookup_cache_total{mode="exact"}[10m])), 1e-9)
) < 50
```

## 7) Operasyon Notları
- İlk deploy sonrası 5-10 dk warmup beklenmelidir.
- Trafik düşükse `rate([5m])` yerine `increase([15m])` daha anlamlı olabilir.
- Module bazlı alarm eşikleri trafik profiline göre ayrı belirlenmelidir.
