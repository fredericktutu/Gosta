result_zero=(0 1 0 2)
result_cap=(0 5 120 2)
result_controlflow=(82 82 18 2)
result_pipeline=(137 120 0)


testSuite(){
    result_array=$1
    fileprefix="testcase/$2/case"
    logprefix="log/$2/case"
	weight=$3
    
    echo "---run testsuite '$2'---"
    pass=0
    fail=0
    n=0
	
    for expect in ${result_array[*]}; do
	    ((n=$n+1))
	    filepath="${fileprefix}${n}.go"
		exepath="case${n}.exe"
	    logpath="${logprefix}${n}.txt"
	    # echo "filepath: $filepath , logpath $logpath"
		go build ${filepath}
		if [ "$?" != "0" ]; then 
			echo "case${n}: can't compile!"
			continue
		fi
		rm -rf ${exepath}
		
	    bin/main.exe ${filepath} ${weight}  --exec > ${logpath}
	    result=$?
	    if [ "$expect" == "$result" ]; then
		    echo "case${n}: pass!"
		    ((pass=$pass+1))
            else
		    echo "case${n}: fail! expect ${expect}, find ${result}"
		    ((fail=$fail+1))
        fi
    done
    echo "---conclusion: ${pass} passes, ${fail} fails---"
}


suite=$1  #input test_suite
weight=$2
case $suite in
	--zero) testSuite  "${result_zero[*]}"  "test_zero" "$weight"
	;;
	--cap) testSuite "${result_cap[*]}" "test_cap" "$weight"
	;;
	--pipeline) testSuite "${result_pipeline[*]}" "test_pipeline" "$weight"
	;;
	--controlflow) testSuite "${result_controlflow[*]}" "test_controlflow" "$weight"
	;;
esac
