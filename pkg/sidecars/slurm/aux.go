package slurm

import (
	"log"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
)

func prepare_envs(container v1.Container) []string {
	env := make([]string, 1)
	env = append(env, "--env")
	env_data := ""
	for _, env_var := range container.Env {
		tmp := (env_var.Name + "=" + env_var.Value + ",")
		env_data += tmp
	}
	if last := len(env_data) - 1; last >= 0 && env_data[last] == ',' {
		env_data = env_data[:last]
	}
	env = append(env, env_data)
	return env
}

func prepare_mounts(container v1.Container) []string {
	mount := make([]string, 1)
	mount = append(mount, "--bind")
	mount_data := ""
	pod_name := strings.Split(container.Name, "-")

	err := os.MkdirAll(".knoc/"+strings.Join(pod_name[:6], "-"), os.ModePerm)
	if err != nil {
		log.Fatalln("Cant create directory")
	}
	for _, mount_var := range container.VolumeMounts {
		f, err := os.Create(".knoc/" + strings.Join(pod_name[:6], "-") + "/" + mount_var.Name)
		f.WriteString("")
		if err != nil {
			log.Fatalln("Cant create directory")
		}
		path := (".knoc/" + strings.Join(pod_name[:6], "-") + "/" + mount_var.Name + ":" + mount_var.MountPath + ",")
		mount_data += path
	}
	path_hardcoded := ("/cvmfs/grid.cern.ch/etc/grid-security:/etc/grid-security" + "," +
		"/m100_scratch/userexternal/dspiga00:/m100_scratch/userexternal/dspiga00" + "," +
		"/m100_work:/m100_work" + "," +
		"/cvmfs:/cvmfs" + "," +
		"/m100_work/INF22_lhc_0/CMS/SITECONF:/marconi_work/Pra18_4658/cms/SITECONF" + ",")
	mount_data += path_hardcoded
	if last := len(mount_data) - 1; last >= 0 && mount_data[last] == ',' {
		mount_data = mount_data[:last]
	}
	return append(mount, mount_data)
}
