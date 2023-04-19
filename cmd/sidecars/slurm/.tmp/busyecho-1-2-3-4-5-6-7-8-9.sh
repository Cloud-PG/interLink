#!/bin/bash
#SBATCH --job-name=busyecho-1-2-3-4-5-6-7-8-9
#SBATCH -A
#SBATCH INF23_lhc_0
#SBATCH 
#SBATCH --job-name=TEST_KNOC
#SBATCH -p
#SBATCH m100_usr_prod
#SBATCH -t
#SBATCH 1440
#SBATCH --gres="gpu:0"
#SBATCH --cpus-per-task=1
. ~/.bash_profile
module load singularity
export SINGULARITYENV_SINGULARITY_TMPDIR=$CINECA_SCRATCH
export SINGULARITYENV_SINGULARITY_CACHEDIR=$CINECA_SCRATCH
pwd; hostname; date

singularity exec  --env   --bind .knoc/busyecho-1-2-3-4-5/kube-api-access-hxjkg:/var/run/secrets/kubernetes.io/serviceaccount,/cvmfs/grid.cern.ch/etc/grid-security:/etc/grid-security,/m100_scratch/userexternal/dspiga00:/m100_scratch/userexternal/dspiga00,/m100_work:/m100_work,/cvmfs:/cvmfs,/m100_work/INF23_lhc_0/CMS/SITECONF:/marconi_work/Pra18_4658/cms/SITECONF /cvmfs/unpacked.cern.ch/registry.hub.docker.com/cmssw/el8:ppc64le /m100/home/userexternal/gsurace0/script.sh >> .knoc/busyecho-1-2-3-4-5-6-7-8-9.out 2>> .knoc/busyecho-1-2-3-4-5-6-7-8-9.err 
 echo $? > .knoc/busyecho-1-2-3-4-5-6-7-8-9.status