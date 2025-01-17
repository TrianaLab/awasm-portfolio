document.addEventListener("render-graph", (event) => {
    const globalJsonData = event.detail;

    if (!globalJsonData) {
        console.error("No data available for visualization.");
        return;
    }

    // NEW back arrow icon
    const backButtonSvg = `
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
            <!--!Font Awesome Free 6.7.2 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free -->
            <path d="M125.7 160l50.3 0c17.7 0 32 14.3 32 32s-14.3 32-32 32L48 224c-17.7 0-32-14.3-32-32L16 64c0-17.7 14.3-32 32-32s32 14.3 32 32l0 51.2L97.6 97.6c87.5-87.5 229.3-87.5 316.8 0s87.5 229.3 0 316.8s-229.3 87.5-316.8 0c-12.5-12.5-12.5-32.8 0-45.3s32.8-12.5 45.3 0c62.5 62.5 163.8 62.5 226.3 0s62.5-163.8 0-226.3s-163.8-62.5-226.3 0L125.7 160z"/>
        </svg>
    `;

    /** Utility to compute distance between two points */
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
        // Map the data to bubble-friendly format
        const data = globalJsonData.map((d) => {
            const textLength = d.Name.length;
            const estimatedRadius = textLength * 5 + 20;
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
            .html("") // Remove old content
            .append("svg")
            .attr("width", width)
            .attr("height", height);

        // We'll need the raw node for bounding rect
        const svgNode = svgSelection.node();
        const svgRect = svgNode.getBoundingClientRect();

        // Use D3 color scale for bubble fill
        const colorScale = d3.scaleOrdinal(d3.schemeCategory10);

        // Force Simulation
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

        // Each node is a group <g>
        const node = svgSelection
            .selectAll("g")
            .data(data)
            .enter()
            .append("g")
            .attr("class", "node");

        // Bubbles (circles)
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

        // Labels
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

        /**
         * Zoom into the bubble with a circle clip-path,
         * show details in YAML (using js-yaml),
         * and provide a back arrow to close.
         * 
         * Font: "Courier New" for the overlay content
         */
        function zoomIntoBubble(d) {
            // Center of this bubble in absolute coords
            const bubbleCenterX = svgRect.left + d.x;
            const bubbleCenterY = svgRect.top + d.y;

            // Compute how big the circle must grow to cover the screen
            const screenWidth = window.innerWidth;
            const screenHeight = window.innerHeight;

            const radiusNeeded = Math.max(
                distance(bubbleCenterX, bubbleCenterY, 0, 0),
                distance(bubbleCenterX, bubbleCenterY, screenWidth, 0),
                distance(bubbleCenterX, bubbleCenterY, 0, screenHeight),
                distance(bubbleCenterX, bubbleCenterY, screenWidth, screenHeight)
            );

            // 1) Create an overlay that uses a circular clip-path
            const overlay = d3.select("body")
                .append("div")
                .attr("class", "bubble-overlay")
                .style("clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                .style("-webkit-clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`);

            // 2) Circle clip expand animation (zoom in)
            overlay
                .transition()
                .duration(500)
                .style("clip-path", `circle(${radiusNeeded}px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                .style("-webkit-clip-path", `circle(${radiusNeeded}px at ${bubbleCenterX}px ${bubbleCenterY}px)`);

            // 3) Add a back arrow at the top-left
            const backButton = overlay
                .append("div")
                .attr("class", "back-arrow")
                .html(backButtonSvg);

            // 4) Main content container for the details,
            //    using "Courier New" or fallback
            const contentContainer = overlay
                .append("div")
                .attr("class", "bubble-content")
                .style("font-family", "'Courier New', Courier, monospace");

            // Resource name
            contentContainer
                .append("h1")
                .style("margin-bottom", "0.5em")
                .text(d.name);

            // YAML details block (using js-yaml)
            const yamlStr = jsyaml.dump(d.details);

            contentContainer
                .append("pre")
                .style("text-align", "left")
                .style("white-space", "pre-wrap")
                .style("word-break", "break-word")
                .text(yamlStr);

            // 5) Click arrow → circle clip collapse animation (zoom out)
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

        // Start force simulation
        simulation.alpha(1).restart();
    }

    ensureUICanvasVisible();
});
