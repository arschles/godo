package actions

import (
	"github.com/arschles/gci/build"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
)

func Lint(c *cli.Context) {
	bfile := build.GetFileOrDie(c.GlobalString(FlagFile))

	if bfile.Version > BuildFileVersion {
		log.Err("The build file has a higher version (%d) than this build supports (<= %d)", bfile.Version, BuildFileVersion)
	}
	for i, v := range bfile.Vars {
		if v.Default == "" {
			log.Err("Var %d (%s) doesn't have a default specified. It needs at least a default", i, v.Name)
		}
	}
	for i, step := range bfile.StepIncludes {
		if step.Name == "" {
			log.Err("Step include %d doesn't have a name", i)
			continue
		}
		if step.Path == "" {
			log.Err("Step include %d (%s) doesn't have a path", i, step.Name)
		}
	}

	for i, pipeline := range bfile.Pipelines {
		if err := pipeline.Validate(); err != nil {
			log.Err("Pipeline %d: %s", i, pipeline.Name, err)
			continue
		}
		for j, step := range pipeline.Steps {
			if err := step.Validate(); err != nil {
				log.Err("Pipeline %d (%s), step %d: %s", i, pipeline.Name, j, step.Name, err)
				continue
			}

			for z, param := range step.Params {
				if err := param.Validate(); err != nil {
					log.Err("Pipeline %d (%s), step %d (%s), param %d: %s", i, pipeline.Name, j, step.Name, z, err)
					continue
				}
			}

		}
	}
}

func lintStep(step build.PipelineStep) {

}
