document.addEventListener("render-graph", (event) => {
    const globalJsonData = event.detail;

    if (!globalJsonData) {
        console.error("No data available for visualization.");
        return;
    }

    // Transform JSON Data for D3.js
    const data = globalJsonData.map((d) => ({
        name: d.Name,
        namespace: d.Namespace || "unknown",
        value: Math.random() * 800 + 500,
        radius: Math.sqrt(Math.random() * 800 + 500),
        zoomed: false,
        details: d,
        x: Math.random() * window.innerWidth, // Random starting x position
        y: Math.random() * window.innerHeight, // Random starting y position
    }));

    function renderGraph(data) {
        const width = window.innerWidth;
        const height = window.innerHeight;

        const svg = d3.select("#graph-container")
            .html("") // Clear existing content
            .append("svg")
            .attr("width", width)
            .attr("height", height)
            .style("background-color", "#1e1e1e");

        const colorScale = d3.scaleOrdinal(d3.schemeCategory10);

        const simulation = d3.forceSimulation(data)
            .force("center", d3.forceCenter(width / 2, height / 2).strength(0.5)) // Strong pull to the center
            .force("collision", d3.forceCollide().radius((d) => d.radius).strength(1)) // Proper collision handling
            .force("x", d3.forceX(width / 2).strength(0.05)) // Gentle horizontal pull
            .force("y", d3.forceY(height / 2).strength(0.05)) // Gentle vertical pull
            .on("tick", ticked);

        const node = svg
            .selectAll("g")
            .data(data)
            .enter()
            .append("g")
            .attr("class", "node");

        // Append Circles
        node.append("circle")
            .attr("r", (d) => d.radius)
            .attr("fill", (d) => colorScale(d.namespace));

        // Append Text
        node.append("text")
            .text((d) => d.name)
            .attr("dy", "0.3em")
            .attr("text-anchor", "middle")
            .style("fill", "#ffffff")
            .style("font-size", (d) => `${d.radius / 6}px`);

        function ticked() {
            node.attr("transform", (d) => `translate(${d.x}, ${d.y})`);
        }

        simulation.alpha(1).restart(); // Reset alpha to ensure forces are applied
    }

    renderGraph(data);
});
