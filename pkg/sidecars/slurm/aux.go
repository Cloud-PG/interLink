package slurm

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	exec2 "github.com/alexellis/go-execute/pkg/v1"
	commonIL "github.com/cloud-pg/interlink/pkg/common"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JidStruct struct {
	JID string
	Pod v1.Pod
}

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

	if len(pod_name) > 6 {
		pod_name = pod_name[0:6]
	}

	err := os.MkdirAll(".knoc/"+strings.Join(pod_name[:len(pod_name)-1], "-"), os.ModePerm)
	if err != nil {
		log.Fatalln("Cant create directory")
	}

	for _, mount_var := range container.VolumeMounts {

		f, err := os.Create(".knoc/" + strings.Join(pod_name[:len(pod_name)-1], "-") + "/" + mount_var.Name)
		f.WriteString("")
		if err != nil {
			log.Fatalln("Cant create directory")
		}

		path := (".knoc/" + strings.Join(pod_name[:len(pod_name)-1], "-") + "/" + mount_var.Name + ":" + mount_var.MountPath + ",")
		mount_data += path
	}
	path_hardcoded := ("/cvmfs/grid.cern.ch/etc/grid-security:/etc/grid-security" + "," +
		"/m100_scratch/userexternal/dspiga00:/m100_scratch/userexternal/dspiga00" + "," +
		"/m100_work:/m100_work" + "," +
		"/cvmfs:/cvmfs" + "," +
		"/m100_work/INF23_lhc_0/CMS/SITECONF:/marconi_work/Pra18_4658/cms/SITECONF" + ",")
	mount_data += path_hardcoded
	if last := len(mount_data) - 1; last >= 0 && mount_data[last] == ',' {
		mount_data = mount_data[:last]
	}
	return append(mount, mount_data)
}

func produce_slurm_script(container v1.Container, metadata metav1.ObjectMeta, command []string) string {
	newpath := filepath.Join(".", ".tmp")
	err := os.MkdirAll(newpath, os.ModePerm)
	f, err := os.Create(".tmp/" + container.Name + ".sh")
	if err != nil {
		log.Fatalln("Cant create slurm_script")
	}
	var sbatch_flags_from_argo []string
	var sbatch_flags_as_string = ""
	if slurm_flags, ok := metadata.Annotations["slurm-job.knoc.io/flags"]; ok {
		sbatch_flags_from_argo = strings.Split(slurm_flags, " ")
		log.Print(sbatch_flags_from_argo)
	}
	if mpi_flags, ok := metadata.Annotations["slurm-job.knoc.io/mpi-flags"]; ok {
		if mpi_flags != "true" {
			mpi := append([]string{"mpiexec", "-np", "$SLURM_NTASKS"}, strings.Split(mpi_flags, " ")...)
			command = append(mpi, command...)
		}
		log.Print(mpi_flags)
	}
	for _, slurm_flag := range sbatch_flags_from_argo {
		sbatch_flags_as_string += "\n#SBATCH " + slurm_flag
	}

	prefix := ""
	postfix := ""

	if commonIL.InterLinkConfigInst.Tsocks {
		postfix += "\n\nkill -15 $SSH_PID &> log2.txt"

		prefix += "\nmin_port=10000"
		prefix += "\nmax_port=65000"
		prefix += "\nfor ((port=$min_port; port<=$max_port; port++))"
		prefix += "\ndo"
		prefix += "\n  temp=$(ss -tulpn | grep :$port)"
		prefix += "\n  if [ -z \"$temp\" ]"
		prefix += "\n  then"
		prefix += "\n    break"
		prefix += "\n  fi"
		prefix += "\ndone"

		prefix += "\nssh -4 -N -D $port " + commonIL.InterLinkConfigInst.Tsockslogin + " &"
		prefix += "\nSSH_PID=$!"
		prefix += "\nexport LD_PRELOAD=" + commonIL.InterLinkConfigInst.Tsockspath
	}

	if commonIL.InterLinkConfigInst.Commandprefix != "" {
		prefix += "\n" + commonIL.InterLinkConfigInst.Commandprefix
	}

	sbatch_macros := "#!/bin/bash" +
		"\n#SBATCH --job-name=" + container.Name +
		sbatch_flags_as_string +
		"\n. ~/.bash_profile" +
		"\nmodule load singularity" +
		"\nexport SINGULARITYENV_SINGULARITY_TMPDIR=$CINECA_SCRATCH" +
		"\nexport SINGULARITYENV_SINGULARITY_CACHEDIR=$CINECA_SCRATCH" +
		"\npwd; hostname; date" +
		prefix +
		"\n"
	f.WriteString(sbatch_macros + "\n" + strings.Join(command[:], " ") + " >> " + ".knoc/" + container.Name + ".out 2>> " + ".knoc/" + container.Name + ".err \n echo $? > " + ".knoc/" + container.Name + ".status" + postfix)
	f.Close()
	return ".tmp/" + container.Name + ".sh"
}

func slurm_batch_submit(path string) string {
	cmd := []string{path}
	shell := exec2.ExecTask{
		Command: "sbatch",
		Args:    cmd,
		Shell:   true,
	}

	execReturn, _ := shell.Execute()
	execReturn.Stdout = strings.ReplaceAll(execReturn.Stdout, "\n", "")

	if execReturn.Stderr != "" {
		log.Println("Could not run sbatch. " + execReturn.Stderr)
	}
	return string(execReturn.Stdout)
}

func handle_jid(container v1.Container, output string, pod v1.Pod) {
	r := regexp.MustCompile(`Submitted batch job (?P<jid>\d+)`)
	jid := r.FindStringSubmatch(output)
	f, err := os.Create(".knoc/" + container.Name + ".jid")
	if err != nil {
		log.Println("Cant create jid_file")
	}
	f.WriteString(jid[1])
	JID = append(JID, JidStruct{JID: jid[1], Pod: pod})
	f.Close()
}

func delete_container(container v1.Container) {
	data, err := os.ReadFile(".knoc/" + container.Name + ".jid")
	if err != nil {
		log.Fatalln("Can't find job id of container")
	}
	jid, err := strconv.Atoi(string(data))
	if err != nil {
		log.Fatalln("Can't find job id of container")
	}
	_, err = exec.Command(commonIL.SCANCEL, fmt.Sprint(jid)).Output()
	if err != nil {
		log.Println("Could not delete job", jid)
	} else {
		log.Println("Successfully deleted job ", jid)
	}
	exec.Command("rm", "-f ", ".knoc/"+container.Name+".out")
	exec.Command("rm", "-f ", ".knoc/"+container.Name+".err")
	exec.Command("rm", "-f ", ".knoc/"+container.Name+".status")
	exec.Command("rm", "-f ", ".knoc/"+container.Name+".jid")
	exec.Command("rm", "-rf", " .knoc/"+container.Name)
}
