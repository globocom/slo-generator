# slo-generator
Easily setup a service level objective using prometheus, based on lessons from the book linked below: https://landing.google.com/sre/workbook/chapters/alerting-on-slos/.


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

# Grafana integration

All generated SLOs are visible by grafana:

![Overview](https://github.com/globocom/slo-generator/raw/master/grafana-screenshots/slo-overview.png)
![Long Term](https://github.com/globocom/slo-generator/raw/master/grafana-screenshots/slo-long-term.png)

Import dashboard using following [JSON files](./grafana-dashboards)
