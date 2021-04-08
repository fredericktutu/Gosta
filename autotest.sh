result_zero=(0 1 0 2)
result_cap=()
result_pipeline=()
result_controlflow=()

testSuite(){
    result_array=$1
    fileprefix="testcase/$2/case"
    logprefix="log/$2/case"
    
    echo "---run testsuite '$2'---"
    pass=0
    fail=0
    n=0
    for expect in ${result_array[*]}; do
	    ((n=$n+1))
	    filepath="${fileprefix}${n}.go"
	    logpath="${logprefix}${n}.txt"
	    # echo "filepath: $filepath , logpath $logpath"
	    bin/main.exe ${filepath} --exec > ${logpath}
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
case $suite in
	--zero) testSuite  "${result_zero[*]}" "test_zero"
	;;
	--cap) testSuite "${result_cap[*]}" "test_cap"
	;;
	--pipeline) testSuite "${result_pipeline[*]}" "test_pipeline"
	;;
	--controlflow) testSuite "${result_controlflow[*]}" "test_controlflow"
	;;
esac
