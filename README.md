# slo-generator
Easily setup a service level objective using prometheus, based on lessons from the [SRE workbook](https://landing.google.com/sre/workbook/chapters/alerting-on-slos/).

```
slo-generator -slo.path=slo_example.yml -rule.output rule.yml
```


# Generating

Look the file [slo_example.yml](./examples/slo_example.yml) to see how to parametrize SLOs and generate Prometheus rules by running the following command:

# Alert methods currently supported

- [x] 1. Target Error Rate ≥ SLO Threshold, using `alertMethod: simple`
- [x] 2. Increased Alert Window, using `alertMethod: simple`
- [x] 3. Incrementing Alert Duration, using `alertMethod: simple`
- [x] 4. Alert on Burn Rate `alertMethod: simple and burnRate: <rate>`
- [x] 5. Multiple Burn Rate Alerts, using `alertMethod: multi-window and shortWindow: false`
- [x] 6. Multiwindow, Multi-Burn-Rate Alerts, using `alertMethod: multi-window`

## alertMethod: simple

Look the file [slo_simple_example.yml](./examples/slo_simple_example.yml) to see a full example of usage.
the simple alert method require two params:

1. alertWindow: how far back in time will used to alerting. supported values: 5m, 30m, 1h, 2h, 6h, 1d and 3d.
2. alertWait: for long time will begin fire an alert.

## alertMethod: multi-window

The philosofy of this alert is described on the section of book: (https://landing.google.com/sre/workbook/chapters/alerting-on-slos#6-multiwindow-multi-burn-rate-alerts)

# SLOs at scale

The Workbook suggests to create classes to simplify how to set a SLO for your services, read details about concepts [here](https://landing.google.com/sre/workbook/chapters/alerting-on-slos/#alerting_at_scale)

Look at [slo_example_with_classes.yml](./examples/slo_example_with_classes.yml) and [slo_classes.yml](./examples/slo_classes.yml) to see how to define classes and associate with your services.

After that, you can run the command specifying the classes:

```
slo-generator -slo.path=slo_example_with_classes.yml -classes.path slo_classes.yml -rule.output rule.yml
```

# Kubernetes integration

We support to export SLOs as a well known [PrometheusRule](https://github.com/prometheus-operator/prometheus-operator) resources managed by prometheus-operator, just use `-kubernetes` flag, example:

```
slo-generator -kubernetes -slo.path=slo_example.yml > slo_manifest.yml
cat slo_manifest.yml
kubectl apply -f slo_manifest.yml
```

# Grafana integration

All generated SLOs are visible by grafana:

![Overview](https://github.com/globocom/slo-generator/raw/master/grafana-screenshots/slo-overview.png)
![Long Term](https://github.com/globocom/slo-generator/raw/master/grafana-screenshots/slo-long-term.png)

Import dashboard using following [JSON files](./grafana-dashboards)
