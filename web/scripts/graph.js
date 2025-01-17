document.addEventListener("render-graph", (event) => {
    const globalJsonData = event.detail;

    if (!globalJsonData) {
        console.error("No data available for visualization.");
        return;
    }

    function ensureUICanvasVisible(attempt = 1, maxAttempts = 10) {
        const uiCanvas = document.querySelector("#ui-canvas");
        const canvasRect = uiCanvas.getBoundingClientRect();
        const parentRect = uiCanvas.parentElement.getBoundingClientRect();

        if (
            canvasRect.width === 0 ||
            canvasRect.height === 0 ||
            parentRect.width === 0 ||
            parentRect.height === 0
        ) {
            console.warn(`[Warning] #ui-canvas or parent has zero dimensions. Retrying... (${attempt}/${maxAttempts})`);

            if (attempt === 1) {
                uiCanvas.style.display = "flex";
                uiCanvas.style.width = "100vw";
                uiCanvas.style.height = "100vh";
                uiCanvas.style.visibility = "visible";
                uiCanvas.style.opacity = "1";
            }

            if (attempt < maxAttempts) {
                setTimeout(() => ensureUICanvasVisible(attempt + 1, maxAttempts), 100);
            } else {
                console.error("#ui-canvas still has zero dimensions after 10 retries.");
            }
        } else {
            console.info("#ui-canvas is now visible with valid dimensions.");
            renderGraphContent(); // Proceed with graph rendering
        }
    }

    function renderGraphContent() {
        const data = globalJsonData.map((d) => {
            const textLength = d.Name.length;
            const estimatedRadius = textLength * 5 + 20; // Smaller default radius
            return {
                name: d.Name,
                namespace: d.Namespace || "unknown",
                radius: estimatedRadius,
                details: d,
                x: Math.random() * window.innerWidth,
                y: Math.random() * window.innerHeight,
            };
        });

        console.log("[Data Mapped] Final Data for Graph Rendering:", data);

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
            .force("center", d3.forceCenter(width / 2, height / 2).strength(0.3))
            .force("collision", d3.forceCollide().radius((d) => d.radius).strength(1))
            .force("x", d3.forceX(width / 2).strength(0.05))
            .force("y", d3.forceY(height / 2).strength(0.05))
            .on("tick", ticked);

        const node = svg
            .selectAll("g")
            .data(data)
            .enter()
            .append("g")
            .attr("class", "node");

        // Append Circles
        const circles = node.append("circle")
            .attr("r", (d) => d.radius)
            .attr("fill", (d) => colorScale(d.namespace))
            .attr("stroke", "#ffffff")
            .attr("stroke-width", 1)
            .style("opacity", 1) // Fully opaque
            .on("mouseover", function (event, d) {
                d3.select(this)
                    .transition()
                    .duration(200)
                    .attr("r", d.radius * 1.5); // Zoom only this bubble

                simulation.force(
                    "collision",
                    d3.forceCollide().radius((dInner) => {
                        return dInner === d ? d.radius * 1.5 : dInner.radius;
                    })
                );
                simulation.alpha(0.8).restart(); // Restart simulation to adapt to the new size
            })
            .on("mouseout", function (event, d) {
                d3.select(this)
                    .transition()
                    .duration(200)
                    .attr("r", d.radius); // Reset size

                simulation.force(
                    "collision",
                    d3.forceCollide().radius((dInner) => dInner.radius)
                );
                simulation.alpha(0.8).restart(); // Restart simulation
            })
            .on("click", (event, d) => {
                zoomIntoBubble(d);
            });

        // Append Text
        node.append("text")
            .text((d) => d.name)
            .attr("dy", "0.3em")
            .attr("text-anchor", "middle")
            .style("fill", "#ffffff")
            .style("font-family", "Courier, monospace")
            .style("font-size", "14px")
            .style("pointer-events", "none");

        function ticked() {
            node.attr("transform", (d) => `translate(${d.x}, ${d.y})`);
        }

        function zoomIntoBubble(d) {
            const overlay = d3.select("body").append("div")
                .attr("class", "bubble-overlay")
                .style("position", "fixed")
                .style("top", 0)
                .style("left", 0)
                .style("width", "100vw")
                .style("height", "100vh")
                .style("background-color", colorScale(d.namespace))
                .style("display", "flex")
                .style("justify-content", "center")
                .style("align-items", "center")
                .style("z-index", 1000)
                .style("opacity", 0)
                .transition()
                .duration(500)
                .style("opacity", 1);

            overlay.append("div")
                .style("color", "#ffffff")
                .style("font-family", "Courier, monospace")
                .style("text-align", "center")
                .style("max-width", "80%")
                .style("padding", "20px")
                .html(`
                    <h1>${d.name}</h1>
                    <p>Namespace: ${d.namespace}</p>
                    <p>Details:</p>
                    <pre>${JSON.stringify(d.details, null, 2)}</pre>
                    <button id="close-overlay" style="margin-top: 20px; padding: 10px 20px; font-size: 16px;">Close</button>
                `);

            d3.select("#close-overlay").on("click", () => {
                overlay.transition()
                    .duration(500)
                    .style("opacity", 0)
                    .remove();
            });
        }

        simulation.alpha(1).restart();
    }

    ensureUICanvasVisible();
});
