package actions

import (
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
)

// runTarget runs a target's dependencies in order, and then runs name. it prints out information about the target's execution, or calls log.Die to notify of errors, such as a runtime error, missing dependency, or a dependency cycle
// func runTarget(buildFile *build.File, targets map[string]config.Target, target build.Target, visited map[string]struct{}) {
// log.Die("TODO")
// if _, ok := visited[target.Name]; ok {
// 	log.Die("target %s already has a dependency", target.Name)
// 	return
// }
// if target.Depends == "" {
// 	log.Info("target %s", target.Name)
// 	for _, cmd := range target.Commands {
// 		if cmd == "" {
// 			log.Die("command is empty for target %s", target.Name)
// 		}
// 		cmdSpl := strings.Split(cmd, " ")
// 		cmd := exec.Command(cmdSpl[0], cmdSpl[1:]...)
// 		out := runOrDie(cmd, append(os.Environ(), Envs(consfile.Envs).Strings()...))
// 		if len(out) > 0 {
// 			log.Info(string(out))
// 		}
// 	}
// 	return
// }
//
// dependencyTarget, ok := targets[target.Depends]
// if !ok {
// 	log.Die("target %s has dependency %s, which doesn't exist", target.Name, target.Depends)
// 	return
// }
// visited[target.Name] = struct{}{}
// runTarget(consfile, targets, dependencyTarget, visited)
// }

func Run(c *cli.Context) {
	log.Die("TODO")
	// consfile := getConsfileOrDie()
	// tgtName := c.Args().First()
	// if tgtName == "" {
	// 	log.Die("no target given")
	// }
	// targetsMap := make(map[string]Target)
	// var tgt *Target = nil
	// for _, target := range consfile.Targets {
	// 	if _, ok := targetsMap[target.Name]; ok {
	// 		log.Die("target name %s is duplicated", target.Name)
	// 		return
	// 	}
	// 	targetsMap[target.Name] = target
	// 	if target.Name == tgtName {
	// 		tgt = &target
	// 	}
	// }
	// if tgt == nil {
	// 	log.Die("no target %s", tgtName)
	// 	return
	// }
	// runTarget(consfile, targetsMap, *tgt, make(map[string]struct{}))
	// log.Info("done")
}
