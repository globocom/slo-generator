# slo-generator
Easy setup a service level objective using prometheus, based on lessons of book: https://landing.google.com/sre/workbook/chapters/alerting-on-slos/.


# Alert methods currently supported

- [ ] 1. Target Error Rate ≥ SLO Threshold
- [ ] 2. Increased Alert Window
- [ ] 3. Incrementing Alert Duration
- [ ] 4. Alert on Burn Rate
- [ ] 5. Multiple Burn Rate Alerts
- [x] 6. Multiwindow, Multi-Burn-Rate Alerts

# Generating

Lookup [slo_example.yml](./slo_example.yml) to parametize SLO and generate prometheus rules running

```
slo-generator -slo.path=slo_example.yml -rule.output rule.yml
```
