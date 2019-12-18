# slo-generator
Easily setup a service level objective using prometheus, based on lessons from the [SRE workbook](https://landing.google.com/sre/workbook/chapters/alerting-on-slos/).


# Alert methods currently supported

- [ ] 1. Target Error Rate ≥ SLO Threshold
- [ ] 2. Increased Alert Window
- [ ] 3. Incrementing Alert Duration
- [ ] 4. Alert on Burn Rate
- [ ] 5. Multiple Burn Rate Alerts
- [x] 6. Multiwindow, Multi-Burn-Rate Alerts

# Generating

Lookup [slo_example.yml](./examples/slo_example.yml) to parametize SLO and generate prometheus rules running

```
slo-generator -slo.path=slo_example.yml -rule.output rule.yml
```

# SLOs at scale

The Workbook suggests to create classes to simplify how to set a SLO for your services, read details about concepts [here](https://landing.google.com/sre/workbook/chapters/alerting-on-slos/#alerting_at_scale)

Lookup [slo_example_with_classes.yml](./examples/slo_example_with_classes.yml) and [slo_classes.yml](./examples/slo_classes.yml) to see how to define classes and associate with your services.

After that, you can run command specifing the classes file like this:

```
slo-generator -slo.path=slo_example_with_classes.yml -classes.path slo_classes.yml -rule.output rule.yml
```

# Grafana integration

All generated SLOs are visible by grafana:

![Overview](https://github.com/globocom/slo-generator/raw/master/grafana-screenshots/slo-overview.png)
![Long Term](https://github.com/globocom/slo-generator/raw/master/grafana-screenshots/slo-long-term.png)

Import dashboard using following [JSON files](./grafana-dashboards)
