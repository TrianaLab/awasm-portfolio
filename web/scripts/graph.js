document.addEventListener("render-graph", (event) => {
    const globalJsonData = event.detail;

    if (!globalJsonData) {
        console.error("No data available for visualization.");
        return;
    }

    const backButtonSvg = `
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512">
            <path d="M9.4 233.4c-12.5 12.5-12.5 32.8 0 45.3l160 160c12.5 12.5 32.8 12.5 45.3 0s12.5-32.8 0-45.3L109.2 288 416 288c17.7 0 32-14.3 32-32s-14.3-32-32-32l-306.7 0L214.6 118.6c12.5-12.5 12.5-32.8 0-45.3s-32.8-12.5-45.3 0l-160 160z"/>
        </svg>
    `;

    let simulation;

    function renderGraphContent() {
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
    
        const svgSelection = d3
            .select("#graph-container")
            .html("")
            .append("svg")
            .attr("width", width)
            .attr("height", height);
    
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
    
        node.append("circle")
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
            })
            .on("click", function (event, d) {
                zoomIntoBubble(d);
            });
    
        node.append("text")
            .text((d) => d.name)
            .attr("dy", "0.3em")
            .attr("text-anchor", "middle")
            .style("fill", "#ffffff")
            .style("font-family", "Courier, monospace")
            .style("font-size", "14px");
    
        function ticked() {
            node.attr("transform", (d) => `translate(${d.x}, ${d.y})`);
        }    

        function zoomIntoBubble(d) {
            const bubbleCenterX = d.x;
            const bubbleCenterY = d.y;
        
            const overlay = d3.select("body")
                .append("div")
                .attr("class", "bubble-overlay")
                .style("clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                .style("-webkit-clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`);
        
            overlay
                .transition()
                .duration(500)
                .style("clip-path", `circle(150% at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                .style("-webkit-clip-path", `circle(150% at ${bubbleCenterX}px ${bubbleCenterY}px)`);
        
            const headerContainer = overlay.append("div")
                .attr("class", "header-container")
                .style("display", "flex")
                .style("align-items", "center") // Vertically center the arrow and header
                .style("gap", "10px") // Add spacing between arrow and header
                .style("padding-left", "10px"); // Align with details padding
        
            headerContainer.append("div")
                .attr("class", "back-arrow")
                .html(backButtonSvg)
                .style("cursor", "pointer")
                .on("click", () => {
                    overlay
                        .transition()
                        .duration(500)
                        .style("clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                        .style("-webkit-clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                        .on("end", () => {
                            overlay.remove();
                        });
                });
        
            headerContainer.append("h1")
                .text(d.name)
                .style("margin", "0") // Remove default margin
                .style("font-family", "Courier, monospace")
                .style("font-size", "24px");
        
            const contentContainer = overlay.append("div")
                .attr("class", "bubble-content")
                .style("padding", "10px"); // Align with header
        
            const yamlStr = jsyaml.dump(d.details);
        
            contentContainer.append("pre")
                .style("text-align", "left")
                .style("white-space", "pre-wrap")
                .style("word-break", "break-word")
                .style("font-family", "Courier, monospace")
                .text(yamlStr);
        }
        
    }

    function updateGraphSize() {
        const svg = d3.select("#graph-container svg");
        if (!svg.empty() && simulation) {
            const width = window.innerWidth;
            const height = window.innerHeight;

            svg.attr("width", width).attr("height", height);

            simulation
                .force("center", d3.forceCenter(width / 2, height / 2))
                .alpha(0.8)
                .restart();
        }
    }

    window.addEventListener("resize", updateGraphSize);

    renderGraphContent();
    updateGraphSize(); // Ensure the initial size is correct
});
