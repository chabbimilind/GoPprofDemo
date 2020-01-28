# GoPprofDemo
Demo of low precision and accuracy of Go pprof

goroutine.go is a parallel code and serial.go is a serial code.

run.sh automates running these two programs.

You can run manually in the following way:

## Parallel program:
goroutine.go has ten exactly the same goroutines: f1-f10 and I use pprof to collect the CPU profiles via. the default OS timers. 
Run the following command several times and notice that each time, pprof reports a different amount of time spent in each of the ten routines.

`go run goroutine.go &&  go tool pprof -top goroutine_prof`

Expectation: each function (f1-10) should be attributed with exactly (or almost exactly) 10% of the execution time on each run.


## Serial program:
serial.go has ten functions (A_expect_1_82 - J_expect_18_18).  The function A_expect_1_82 is expected to consume 1.82% of the total execution time and J_expect_18_18 is expected to consume 18.18% of the execution time, and so on. The code is serial and there is complete data dependence between each function and each iteration of the loop in the functions to avoid any hardware-level optimizations. 
Run the following command several times.

`go run serial.go &&  go tool pprof -top serial_prof`

Expectation: the time attribution should roughly follow the distribution shown below in each run.
```

FUNCTION NAME			EXPECTED RELATIVE TIME

A_expect_1_82			1.82%
B_expect_3_64			3.64%
C_expect_5_46			5.46%
D_expect_7_27			7.27%
E_expect_9_09			9.09%
F_expect_10_91			10.91%
G_expect_12_73			12.73%
H_expect_14_546			14.546%
I_expect_16_36			16.36%
J_expect_18_18			18.18%
```

## Details
The experiments demonstrate that pprof CPU profiles lack accuracy (closeness to the ground truth) and precision (repeatability across different runs). The tests are highly predictive in nature and involve little runtime overheads (allocation, GC, system call, etc.). They are carefully designed so that we can compare pprof CPU profiles against our expectations. The evaluation shows a wide difference from expectation. One should not confuse this to be a runtime issue; the issue is with the use of OS timers used for sampling; OS timers are coarse-grained and have a high skid.

## Observation
Run #1, #2, and #3 respectively show the pprof -top output of `goroutine.go`.
You will notice in a single run (say Run #1), f1-f10 have a wide variance in time attributed to them; the expectation is that each of them gets 10% execution time. There is up to 6x difference in the function with the highest execution time (main.f7, 4210ms, in run1) vs. the function with the lowest execution time (main.f9, 700ms, in run1).
This shows poor accuracy of pprof timer-based profiles.

Furthermore, the time attributed to a function widely varies from run to run. Notice how the top-10 ordering changes. The expectation is that the measurements remain the same from run to run. This shows a poor precision of pprof timer-based profiles.

### Run 1:
```
File: goroutine
Type: cpu
Time: Jan 27, 2020 at 3:45pm (PST)
Duration: 6.70s, Total samples = 18060ms (269.37%)
Showing nodes accounting for 18060ms, 100% of 18060ms total
      flat  flat%   sum%        cum   cum%
    4210ms 23.31% 23.31%     4210ms 23.31%  main.f7
    2610ms 14.45% 37.76%     2610ms 14.45%  main.f2
    2010ms 11.13% 48.89%     2010ms 11.13%  main.f6
    1810ms 10.02% 58.91%     1810ms 10.02%  main.f10
    1780ms  9.86% 68.77%     1780ms  9.86%  main.f3
    1410ms  7.81% 76.58%     1410ms  7.81%  main.f1
    1310ms  7.25% 83.83%     1310ms  7.25%  main.f4
    1110ms  6.15% 89.98%     1110ms  6.15%  main.f5
    1110ms  6.15% 96.12%     1110ms  6.15%  main.f8
     700ms  3.88%   100%      700ms  3.88%  main.f9
```

### Run 2:
```
File: goroutine
Type: cpu
Time: Jan 27, 2020 at 3:45pm (PST)
Duration: 6.71s, Total samples = 17400ms (259.39%)
Showing nodes accounting for 17400ms, 100% of 17400ms total
      flat  flat%   sum%        cum   cum%
    3250ms 18.68% 18.68%     3250ms 18.68%  main.f2
    2180ms 12.53% 31.21%     2180ms 12.53%  main.f9
    2100ms 12.07% 43.28%     2100ms 12.07%  main.f1
    1770ms 10.17% 53.45%     1770ms 10.17%  main.f6
    1700ms  9.77% 63.22%     1700ms  9.77%  main.f5
    1550ms  8.91% 72.13%     1550ms  8.91%  main.f4
    1500ms  8.62% 80.75%     1500ms  8.62%  main.f8
    1440ms  8.28% 89.02%     1440ms  8.28%  main.f3
    1390ms  7.99% 97.01%     1390ms  7.99%  main.f10
     520ms  2.99%   100%      520ms  2.99%  main.f7
```


### Run 3:
```
File: goroutine
Type: cpu
Time: Jan 27, 2020 at 3:48pm (PST)
Duration: 6.71s, Total samples = 17.73s (264.31%)
Showing nodes accounting for 17.73s, 100% of 17.73s total
      flat  flat%   sum%        cum   cum%
     3.74s 21.09% 21.09%      3.74s 21.09%  main.f7
     2.08s 11.73% 32.83%      2.08s 11.73%  main.f9
     2.05s 11.56% 44.39%      2.05s 11.56%  main.f2
     1.85s 10.43% 54.82%      1.85s 10.43%  main.f10
     1.78s 10.04% 64.86%      1.78s 10.04%  main.f1
     1.43s  8.07% 72.93%      1.43s  8.07%  main.f3
     1.42s  8.01% 80.94%      1.42s  8.01%  main.f8
     1.18s  6.66% 87.59%      1.18s  6.66%  main.f6
     1.17s  6.60% 94.19%      1.17s  6.60%  main.f5
     1.03s  5.81%   100%      1.03s  5.81%  main.f4
```



The output for `go run serial.go &&  go tool pprof -top serial_prof` for three runs is shown below.
Comparing the flat% (or cum%) against the expected percentage for each function shows a large difference. For example, in Run 31, `main.H_expect_14_546` is expected to have 14.546% execution time, whereas it is attributed 25% execution time. Furthermore, run to run, there is lack of precision, for example `main.I_expect_16_36` is attributed 6.25% (20ms) execution time in Run #1, whereas it is attributed 21.88% (70ms)  execution time in Run#2


### Run 1:
```
File: serial
Type: cpu
Time: Jan 27, 2020 at 1:42pm (PST)
Duration: 501.51ms, Total samples = 320ms (63.81%)
Showing nodes accounting for 320ms, 100% of 320ms total
      flat  flat%   sum%        cum   cum%
      80ms 25.00% 25.00%       80ms 25.00%  main.H_expect_14_546
      80ms 25.00% 50.00%       80ms 25.00%  main.J_expect_18_18
      60ms 18.75% 68.75%       60ms 18.75%  main.G_expect_12_73
      20ms  6.25% 75.00%       20ms  6.25%  main.B_expect_3_64
      20ms  6.25% 81.25%       20ms  6.25%  main.D_expect_7_27
      20ms  6.25% 87.50%       20ms  6.25%  main.F_expect_10_91
      20ms  6.25% 93.75%       20ms  6.25%  main.I_expect_16_36
      10ms  3.12% 96.88%       10ms  3.12%  main.A_expect_1_82
      10ms  3.12%   100%       10ms  3.12%  main.C_expect_5_46
         0     0%   100%      320ms   100%  main.main
         0     0%   100%      320ms   100%  runtime.main
```

### Run 2:
```
File: serial
Type: cpu
Time: Jan 27, 2020 at 1:44pm (PST)
Duration: 501.31ms, Total samples = 320ms (63.83%)
Showing nodes accounting for 320ms, 100% of 320ms total
      flat  flat%   sum%        cum   cum%
      70ms 21.88% 21.88%       70ms 21.88%  main.I_expect_16_36
      50ms 15.62% 37.50%       50ms 15.62%  main.J_expect_18_18
      40ms 12.50% 50.00%       40ms 12.50%  main.E_expect_9_09
      40ms 12.50% 62.50%       40ms 12.50%  main.F_expect_10_91
      40ms 12.50% 75.00%       40ms 12.50%  main.H_expect_14_546
      30ms  9.38% 84.38%       30ms  9.38%  main.D_expect_7_27
      20ms  6.25% 90.62%       20ms  6.25%  main.B_expect_3_64
      20ms  6.25% 96.88%       20ms  6.25%  main.G_expect_12_73
      10ms  3.12%   100%       10ms  3.12%  main.C_expect_5_46
         0     0%   100%      320ms   100%  main.main
         0     0%   100%      320ms   100%  runtime.main
```

### Run 3:
```
File: serial
Type: cpu
Time: Jan 27, 2020 at 1:45pm (PST)
Duration: 501.39ms, Total samples = 310ms (61.83%)
Showing nodes accounting for 310ms, 100% of 310ms total
      flat  flat%   sum%        cum   cum%
     110ms 35.48% 35.48%      110ms 35.48%  main.J_expect_18_18
      70ms 22.58% 58.06%       70ms 22.58%  main.G_expect_12_73
      60ms 19.35% 77.42%       60ms 19.35%  main.F_expect_10_91
      30ms  9.68% 87.10%       30ms  9.68%  main.I_expect_16_36
      20ms  6.45% 93.55%       20ms  6.45%  main.H_expect_14_546
      10ms  3.23% 96.77%       10ms  3.23%  main.B_expect_3_64
      10ms  3.23%   100%       10ms  3.23%  main.C_expect_5_46
         0     0%   100%      310ms   100%  main.main
         0     0%   100%      310ms   100%  runtime.main
```
