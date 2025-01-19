document.addEventListener("render-graph", (event) => {
    const globalJsonData = event.detail;

    if (!globalJsonData) {
        console.error("No data available for visualization.");
        return;
    }

    // Function to update the graph's size dynamically
    const updateGraphSize = () => {
        const svg = d3.select("#graph-container svg");
        if (!svg.empty()) {
            const width = window.innerWidth;
            const height = window.innerHeight;

            svg.attr("width", width).attr("height", height);

            // Update simulation center
            simulation
                .force("center", d3.forceCenter(width / 2, height / 2))
                .alpha(0.8)
                .restart();
        }
    };

    function renderGraphContent() {
        // Map the data to bubble-friendly format
        const data = globalJsonData.map((d) => {
            const textLength = d.Name.length;
            const estimatedRadius = (textLength * 8) / 2 + 10;
            return {
                name: d.Name,
                namespace: d.Namespace || "unknown",
                radius: estimatedRadius,
                details: d,
                x: Math.random() * window.innerWidth,
                y: Math.random() * window.innerHeight,
            };
        });

        const width = window.innerWidth;
        const height = window.innerHeight;

        // Clear or create the SVG
        const svgSelection = d3
            .select("#graph-container")
            .html("")
            .append("svg")
            .attr("width", width)
            .attr("height", height);

        // High-contrast color scale
        const colorScale = d3.scaleOrdinal(d3.schemeDark2);

        simulation = d3
            .forceSimulation(data)
            .force("center", d3.forceCenter(width / 2, height / 2))
            .force("collision", d3.forceCollide().radius((d) => d.radius))
            .force("x", d3.forceX(width / 2))
            .force("y", d3.forceY(height / 2))
            .on("tick", ticked);

        const node = svgSelection
            .selectAll("g")
            .data(data)
            .enter()
            .append("g")
            .attr("class", "node");

        node
            .append("circle")
            .attr("r", (d) => d.radius)
            .attr("fill", (d) => colorScale(d.namespace))
            .attr("stroke", "#ffffff")
            .attr("stroke-width", 1)
            .on("mouseover", function (event, d) {
                d3.select(this)
                    .transition()
                    .duration(200)
                    .attr("r", d.radius * 1.5);
                simulation.force(
                    "collision",
                    d3.forceCollide().radius((inner) => (inner === d ? d.radius * 1.5 : inner.radius))
                );
                simulation.alpha(0.8).restart();
            })
            .on("mouseout", function (event, d) {
                d3.select(this)
                    .transition()
                    .duration(200)
                    .attr("r", d.radius);
                simulation.force(
                    "collision",
                    d3.forceCollide().radius((inner) => inner.radius)
                );
                simulation.alpha(0.8).restart();
            });

        node
            .append("text")
            .text((d) => d.name)
            .attr("dy", "0.3em")
            .attr("text-anchor", "middle")
            .style("fill", "#ffffff")
            .style("font-family", "Courier, monospace")
            .style("font-size", "14px");

        function ticked() {
            node.attr("transform", (d) => `translate(${d.x}, ${d.y})`);
        }
    }

    renderGraphContent();

    // Initialize dynamic resizing for the graph
    window.addEventListener("resize", updateGraphSize);
});
