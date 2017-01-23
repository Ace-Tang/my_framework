#! /bin/bash

set -o pipefail -o errexit

if [ -z "$1" ];then
	echo "please enter sched, exec or all"
	exit
elif [ "$1" != "sched" ] && [ "$1" != "exec" ] && [ "$1" != "all" ];then
	echo "only receive [sched], [exec] and [all]"
	exit
fi

current_dir=$(pwd)
module_name=${current_dir##*/}
gopath="${current_dir}/gocode"
exec_path="${current_dir}/bin"
repo_path="$gopath/src"
#gopath_w="$repo_path/meggy/Godeps/_workspace"

if [ -e "${gopath}" ];then
	rm -rf $gopath
fi
mkdir -p $repo_path

if [ -e "${exec_path}" ];then
	rm -rf $exec_path
fi
mkdir $exec_path

ln -sf ${current_dir} $repo_path
#mv -T "$repo_path/$module_name" "$repo_path/my_framework"
export GOPATH=$gopath

if [ "$1" == "sched" ];then
	build_name="sched"
	rm -f $exec_path/$build_name
	go build -o $exec_path/$build_name $repo_path/$module_name/cmd/sche_app.go
elif [ "$1" == "exec" ];then
	build_name="sim_exec"
	rm -f $exec_path/$build_name
	go build -o $exec_path/$build_name $repo_path/$module_name/cmd/exec_app.go
else
	rm -rf $exec_path
	build_name="sched"
	go build -o $exec_path/$build_name $repo_path/$module_name/cmd/sche_app.go
	build_name="sim_exec"
	go build -o $exec_path/$build_name $repo_path/$module_name/cmd/exec_app.go
fi
