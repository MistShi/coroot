package model

import (
	"github.com/coroot/coroot/timeseries"
	"github.com/coroot/coroot/utils"
	"strings"
)

type AuditReport struct {
	appId        ApplicationId
	ctx          timeseries.Context
	checkConfigs CheckConfigs

	Name    string    `json:"name"`
	Status  Status    `json:"status"`
	Widgets []*Widget `json:"widgets"`
	Checks  []*Check  `json:"checks"`
}

func NewAuditReport(appId ApplicationId, ctx timeseries.Context, checkConfigs CheckConfigs, name string) *AuditReport {
	return &AuditReport{appId: appId, Name: name, ctx: ctx, checkConfigs: checkConfigs}
}

type Widget struct {
	Chart         *Chart         `json:"chart"`
	ChartGroup    *ChartGroup    `json:"chart_group"`
	Table         *Table         `json:"table"`
	LogPatterns   *LogPatterns   `json:"log_patterns"`
	DependencyMap *DependencyMap `json:"dependency_map"`

	Width string `json:"width"`
}

func (c *AuditReport) GetOrCreateChartGroup(title string) *ChartGroup {
	for _, w := range c.Widgets {
		if cg := w.ChartGroup; cg != nil {
			if cg.Title == title {
				return cg
			}
		}
	}
	cg := &ChartGroup{Title: title}
	c.Widgets = append(c.Widgets, &Widget{ChartGroup: cg})
	return cg
}

func (c *AuditReport) GetOrCreateChartInGroup(title string, chartTitle string) *Chart {
	return c.GetOrCreateChartGroup(title).GetOrCreateChart(c.ctx, chartTitle)
}

func (c *AuditReport) GetOrCreateChart(title string) *Chart {
	for _, w := range c.Widgets {
		if ch := w.Chart; ch != nil {
			if ch.Title == title {
				return ch
			}
		}
	}
	ch := NewChart(c.ctx, title)
	c.Widgets = append(c.Widgets, &Widget{Chart: ch})
	return ch
}

func (c *AuditReport) GetOrCreateDependencyMap() *DependencyMap {
	for _, w := range c.Widgets {
		if w.DependencyMap != nil {
			return w.DependencyMap
		}
	}
	dm := &DependencyMap{}
	c.Widgets = append(c.Widgets, &Widget{DependencyMap: dm, Width: "100%"})
	return dm
}

func (c *AuditReport) GetOrCreateTable(header ...string) *Table {
	for _, w := range c.Widgets {
		if t := w.Table; t != nil {
			return t
		}
	}
	t := &Table{Header: header}
	c.Widgets = append(c.Widgets, &Widget{Table: t, Width: "100%"})
	return t
}

func (c *AuditReport) CreateCheck(cfg CheckConfig) *Check {
	ch := &Check{
		Id:                      cfg.Id,
		Title:                   cfg.Title,
		Status:                  OK,
		Unit:                    cfg.Unit,
		ConditionFormatTemplate: cfg.ConditionFormatTemplate,

		typ:             cfg.Type,
		messageTemplate: cfg.MessageTemplate,
		items:           utils.NewStringSet(),
	}
	switch cfg.Id {
	case Checks.SLOAvailability.Id:
		configs := c.checkConfigs.GetAvailability(c.appId)
		if len(configs) > 0 {
			ch.Threshold = configs[0].ObjectivePercentage
		} else {
			ch.Threshold = Checks.SLOAvailability.DefaultThreshold
		}
	case Checks.SLOLatency.Id:
		configs := c.checkConfigs.GetLatency(c.appId)
		if len(configs) > 0 {
			ch.Threshold = configs[0].ObjectivePercentage
			ch.ConditionFormatTemplate = strings.Replace(
				ch.ConditionFormatTemplate,
				"<bucket>",
				FormatLatencyBucket(configs[0].ObjectiveBucket),
				1,
			)
		} else {
			ch.Threshold = Checks.SLOLatency.DefaultThreshold
			ch.ConditionFormatTemplate = strings.Replace(ch.ConditionFormatTemplate, "<bucket>", "100ms", 1)
		}
	default:
		ch.Threshold = c.checkConfigs.GetSimple(cfg.Id, c.appId).Threshold
	}
	c.Checks = append(c.Checks, ch)
	return ch
}
