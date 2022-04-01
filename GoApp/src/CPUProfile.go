...
f, err := os.Create("cpu.prof")
if err != nil {
    log.Fatal(err)
}
err = pprof.StartCPUProfile(f)
if err != nil {
    log.Fatal(err)
}
defer pprof.StopCPUProfile()
...