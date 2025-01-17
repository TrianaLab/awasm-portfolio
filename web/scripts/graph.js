document.addEventListener("render-graph", (event) => {
    const globalJsonData = event.detail;

    if (!globalJsonData) {
        console.error("No data available for visualization.");
        return;
    }

    // The SVG code for the "back arrow" icon
    const backButtonSvg = `
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
            <path d="M459.5 440.6c9.5 7.9 22.8 9.7 34.1 4.4s18.4-16.6 18.4-29l0-320c0-12.4-7.2-23.7-18.4-29s-24.5-3.6-34.1 4.4L288 214.3l0 41.7 0 41.7L459.5 440.6zM256 352l0-96 0-128 0-32c0-12.4-7.2-23.7-18.4-29s-24.5-3.6-34.1 4.4l-192 160C4.2 237.5 0 246.5 0 256s4.2 18.5 11.5 24.6l192 160c9.5 7.9 22.8 9.7 34.1 4.4s18.4-16.6 18.4-29l0-64z"/>
        </svg>
    `;

    /**
     * Utility to compute distance between two points
     * (used to ensure our clip-path radius covers the screen).
     */
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
        // Prepare data for bubble display
        const data = globalJsonData.map((d) => {
            const textLength = d.Name.length;
            const estimatedRadius = textLength * 5 + 20; 
            return {
                name: d.Name,
                namespace: d.Namespace || "unknown",
                radius: estimatedRadius,
                details: d, // The full JSON object
                x: Math.random() * window.innerWidth,
                y: Math.random() * window.innerHeight,
            };
        });

        console.log("[Data Mapped] Final Data for Graph Rendering:", data);

        const width = window.innerWidth;
        const height = window.innerHeight;

        // Create or clear the SVG
        const svgSelection = d3
            .select("#graph-container")
            .html("") // Clear existing content
            .append("svg")
            .attr("width", width)
            .attr("height", height)
            .style("background-color", "#1e1e1e");

        // We'll need the raw node for bounding rect
        const svgNode = svgSelection.node();
        const svgRect = svgNode.getBoundingClientRect();

        // Color scale for bubble backgrounds
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
            .on("click", function (event, d) {
                zoomIntoBubble(d);
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

        // Ticking the simulation
        function ticked() {
            node.attr("transform", (d) => `translate(${d.x}, ${d.y})`);
        }

        /**
         * Zoom into the bubble with a circle clip-path, 
         * show details in YAML (using js-yaml), 
         * and provide a back arrow to close.
         */
        function zoomIntoBubble(d) {
            const bubbleColor = colorScale(d.namespace);

            // Center of this bubble in absolute coords
            // d.x, d.y is relative to the SVG top-left
            const bubbleCenterX = svgRect.left + d.x;
            const bubbleCenterY = svgRect.top + d.y;

            // Compute how big the circle must grow to cover the screen from that center
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
                .style("position", "fixed")
                .style("top", 0)
                .style("left", 0)
                .style("width", "100vw")
                .style("height", "100vh")
                .style("background-color", bubbleColor)
                .style("z-index", 1000)
                // Start with a circle(0) clip
                .style("clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                .style("-webkit-clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                // Ensure we can scroll inside
                .style("overflow", "hidden");

            // 2) Inject dynamic styles for the scrollbar
            const dynamicStyle = document.createElement("style");
            dynamicStyle.innerHTML = `
                .bubble-overlay {
                    scrollbar-width: thin;
                    scrollbar-color: #ffffff88 ${bubbleColor};
                }
                .bubble-overlay::-webkit-scrollbar {
                    width: 8px;
                }
                .bubble-overlay::-webkit-scrollbar-track {
                    background: ${bubbleColor};
                }
                .bubble-overlay::-webkit-scrollbar-thumb {
                    background-color: #ffffff88;
                    border-radius: 4px;
                    border: 1px solid ${bubbleColor};
                }
            `;
            document.head.appendChild(dynamicStyle);

            // 3) Circle clip expand animation (zoom in)
            overlay
                .transition()
                .duration(500)
                .style("clip-path", `circle(${radiusNeeded}px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                .style("-webkit-clip-path", `circle(${radiusNeeded}px at ${bubbleCenterX}px ${bubbleCenterY}px)`);

            // 4) Add a back arrow at the top-left
            const backButton = overlay
                .append("div")
                .style("position", "absolute")
                .style("top", "20px")
                .style("left", "20px")
                .style("width", "32px")
                .style("height", "32px")
                .style("cursor", "pointer")
                .html(backButtonSvg);

            // Color the arrow white (same as text)
            backButton.select("svg").style("fill", "#ffffff");

            // 5) Main content container for the details
            const contentContainer = overlay
                .append("div")
                .style("color", "#ffffff")
                .style("font-family", "Courier, monospace")
                .style("text-align", "center")
                .style("margin", "40px auto")
                .style("max-width", "80%")
                .style("max-height", "80%")    // limit vertical size
                .style("overflow-y", "auto")  // scroll if too tall
                .style("padding", "20px");

            // Resource name
            contentContainer.append("h1")
                .style("margin-bottom", "0.5em")
                .text(d.name);

            // YAML details block (using js-yaml)
            const yamlStr = jsyaml.dump(d.details);

            contentContainer
                .append("pre")
                .style("text-align", "left")
                .style("font-family", "Courier, monospace")
                .style("padding", "1em")
                .style("white-space", "pre-wrap") 
                .style("word-break", "break-word") 
                .text(yamlStr);

            // 6) Back arrow → circle clip collapse animation (zoom out)
            backButton.on("click", () => {
                overlay
                    .transition()
                    .duration(500)
                    .style("clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                    .style("-webkit-clip-path", `circle(0px at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                    .on("end", () => {
                        overlay.remove();
                        dynamicStyle.remove(); // remove injected scrollbar style
                    });
            });
        }

        // Start force simulation
        simulation.alpha(1).restart();
    }

    ensureUICanvasVisible();
});
