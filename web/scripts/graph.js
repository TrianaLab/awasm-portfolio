document.addEventListener("render-graph", (event) => {
    const globalJsonData = event.detail;

    if (!globalJsonData) {
        console.error("No data available for visualization.");
        return;
    }

    // NEW back arrow icon (unchanged)
    const backButtonSvg = `
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512">
            <!--!Font Awesome Free 6.7.2 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2025 Fonticons, Inc.-->
            <path d="M9.4 233.4c-12.5 12.5-12.5 32.8 0 45.3l160 160c12.5 12.5 32.8 12.5 45.3 0s12.5-32.8 0-45.3L109.2 288 416 288c17.7 0 32-14.3 32-32s-14.3-32-32-32l-306.7 0L214.6 118.6c12.5-12.5 12.5-32.8 0-45.3s-32.8-12.5-45.3 0l-160 160z"/>
        </svg>
    `;

    function distance(x1, y1, x2, y2) {
        return Math.sqrt((x2 - x1) ** 2 + (y2 - y1) ** 2);
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
            console.warn(
                `[Warning] #ui-canvas or parent has zero dimensions. Retrying... (${attempt}/${maxAttempts})`
            );

            if (attempt === 1) {
                uiCanvas.style.display = "flex";
                uiCanvas.style.width = "100vw";
                uiCanvas.style.height = "100vh";
                uiCanvas.style.visibility = "visible";
                uiCanvas.style.opacity = "1";
            }

            if (attempt < maxAttempts) {
                setTimeout(
                    () => ensureUICanvasVisible(attempt + 1, maxAttempts),
                    100
                );
            } else {
                console.error(
                    "#ui-canvas still has zero dimensions after 10 retries."
                );
            }
        } else {
            console.info("#ui-canvas is now visible with valid dimensions.");
            renderGraphContent(); // Proceed with graph rendering
        }
    }

    function renderGraphContent() {
        // Map the data to bubble-friendly format with minimal required size
        const data = globalJsonData.map((d) => {
            const textLength = d.Name.length;
            // Estimate radius based on monospace font width (approx 8px per character at 14px font size)
            // Add padding of 10 pixels
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

        console.log("[Data Mapped] Final Data for Graph Rendering:", data);

        const width = window.innerWidth;
        const height = window.innerHeight;

        // Clear or create the SVG
        const svgSelection = d3
            .select("#graph-container")
            .html("")
            .append("svg")
            .attr("width", width)
            .attr("height", height);

        const svgNode = svgSelection.node();
        const svgRect = svgNode.getBoundingClientRect();

        // Use a high-contrast color scale for white text
        const colorScale = d3.scaleOrdinal(d3.schemeDark2);

        const simulation = d3
            .forceSimulation(data)
            .force("center", d3.forceCenter(width / 2, height / 2).strength(0.3))
            .force(
                "collision",
                d3.forceCollide().radius((d) => d.radius).strength(1)
            )
            .force("x", d3.forceX(width / 2).strength(0.05))
            .force("y", d3.forceY(height / 2).strength(0.05))
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
            .style("opacity", 1)
            .on("mouseover", function (event, d) {
                d3.select(this)
                    .transition()
                    .duration(200)
                    .attr("r", d.radius * 1.5);

                simulation.force(
                    "collision",
                    d3.forceCollide().radius((inner) => {
                        return inner === d ? d.radius * 1.5 : inner.radius;
                    })
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
            .on("click", function () {
                zoomIntoBubble(this.__data__);
            });

        node
            .append("text")
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
            const bubbleCenterX = svgRect.left + d.x;
            const bubbleCenterY = svgRect.top + d.y;

            const screenWidth = window.innerWidth;
            const screenHeight = window.innerHeight;

            const radiusNeeded = Math.max(
                distance(bubbleCenterX, bubbleCenterY, 0, 0),
                distance(bubbleCenterX, bubbleCenterY, screenWidth, 0),
                distance(bubbleCenterX, bubbleCenterY, 0, screenHeight),
                distance(bubbleCenterX, bubbleCenterY, screenWidth, screenHeight)
            );

            const overlay = d3.select("body")
                .append("div")
                .attr("class", "bubble-overlay")
                .style("clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                .style("-webkit-clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`);

            overlay
                .transition()
                .duration(500)
                .style("clip-path", `circle(${radiusNeeded}px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                .style("-webkit-clip-path", `circle(${radiusNeeded}px at ${bubbleCenterX}px ${bubbleCenterY}px)`);

            const backButton = overlay
                .append("div")
                .attr("class", "back-arrow")
                .html(backButtonSvg);

            const contentContainer = overlay
                .append("div")
                .attr("class", "bubble-content")
                .style("font-family", "'Courier New', Courier, monospace");

            contentContainer
                .append("h1")
                .style("margin-bottom", "0.5em")
                .text(d.name);

            const yamlStr = jsyaml.dump(d.details);

            contentContainer
                .append("pre")
                .style("text-align", "left")
                .style("white-space", "pre-wrap")
                .style("word-break", "break-word")
                .text(yamlStr);

            backButton.on("click", () => {
                overlay
                    .transition()
                    .duration(500)
                    .style("clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                    .style("-webkit-clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                    .on("end", () => {
                        overlay.remove();
                    });
            });
        }

        simulation.alpha(1).restart();
    }

    ensureUICanvasVisible();
});
