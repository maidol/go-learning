function fibonacci(i) {
	if (i < 2) {
		return i
	}
	return fibonacci(i-2) + fibonacci(i-1)
}

(function(){
  let start = Date.now();
  console.log(fibonacci(50));
  console.log(Date.now() - start)
})()