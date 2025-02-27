package auditor

import (
	"fmt"
	"github.com/coroot/coroot/model"
	"github.com/coroot/coroot/timeseries"
	"github.com/coroot/coroot/utils"
	"github.com/dustin/go-humanize"
	"math"
)

func (a *appAuditor) storage() {
	report := a.addReport("Storage")
	ioCheck := report.CreateCheck(model.Checks.StorageIO)
	spaceCheck := report.CreateCheck(model.Checks.StorageSpace)
	seenVolumes := false
	for _, i := range a.app.Instances {
		for _, v := range i.Volumes {
			fullName := i.Name + ":" + v.MountPoint
			if i.Node != nil {
				seenVolumes = true
				if d := i.Node.Disks[v.Device.Value()]; d != nil {
					report.GetOrCreateChartInGroup("I/O latency <selector>, seconds", v.MountPoint).
						AddSeries(i.Name, d.Await)

					report.GetOrCreateChartInGroup("I/O utilization <selector>, %", v.MountPoint).
						AddSeries(i.Name, d.IOUtilizationPercent)

					if timeseries.Last(d.IOUtilizationPercent) > ioCheck.Threshold {
						ioCheck.AddItem("%s:%s", i.Name, v.MountPoint)
					}

					report.GetOrCreateChartInGroup("IOPS <selector>", fullName).
						Stacked().
						Sorted().
						AddSeries("read", d.ReadOps, "blue").
						AddSeries("write", d.WriteOps, "amber")

					report.GetOrCreateChartInGroup("Bandwidth <selector>, bytes/second", fullName).
						Stacked().
						Sorted().
						AddSeries("read", d.ReadBytes, "blue").
						AddSeries("written", d.WrittenBytes, "amber")

					latencyMs := model.NewTableCell().SetUnit("ms")
					if d.Await != nil {
						latencyMs.SetValue(utils.FormatFloat(timeseries.Last(d.Await) * 1000))
					}
					ioPercent := model.NewTableCell()
					if last := timeseries.Last(d.IOUtilizationPercent); !math.IsNaN(last) {
						ioPercent.SetValue(fmt.Sprintf("%.0f%%", last))
					}
					space := model.NewTableCell()
					capacity := timeseries.Last(v.CapacityBytes)
					usage := timeseries.Last(v.UsedBytes)
					if usage > 0 && capacity > 0 {
						percentage := usage / capacity * 100
						space.SetValue(fmt.Sprintf(
							"%.0f%% (%s / %s)",
							percentage,
							humanize.Bytes(uint64(usage)),
							humanize.Bytes(uint64(capacity))),
						)
						if percentage > spaceCheck.Threshold {
							spaceCheck.AddItem("%s:%s", i.Name, v.MountPoint)
						}
					}
					report.GetOrCreateTable("Volume", "Latency", "I/O", "Space", "Device").AddRow(
						model.NewTableCell(fullName),
						latencyMs,
						ioPercent,
						space,
						model.NewTableCell(v.Device.Value()).AddTag(v.Name.Value()),
					)
				}
				report.GetOrCreateChartInGroup("Disk space <selector>, bytes", fullName).
					Stacked().
					AddSeries("used", v.UsedBytes).
					SetThreshold("total", v.CapacityBytes, timeseries.Max)
			}
		}
	}
	if !seenVolumes {
		ioCheck.SetStatus(model.UNKNOWN, "no volumes found")
		spaceCheck.SetStatus(model.UNKNOWN, "no volumes found")
	}
}
