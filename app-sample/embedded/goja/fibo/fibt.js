function fib(n, a, b) {
    if(n === 0) {
      return a;
    } else if(n === 1) {
      return b;
    }
    return fib(n-1, b, a+b)
}

fib(35, 0, 1);