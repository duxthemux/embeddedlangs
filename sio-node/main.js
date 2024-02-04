let input = [];
process.stdin.on("data", (dt) => {
    input.push(dt)
})
process.stdin.on("end", () => {
    let data = JSON.parse(input)
    data["fld"] = "fromNode;"
    console.log(data)
    process.stdout.write("{}");
})