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
        const yamlData = event.detail;
        const jsonData = typeof yamlData === 'string' ? jsyaml.load(yamlData) : yamlData;

        const groupedData = {};
        jsonData.forEach((item) => {
            const kind = item.Kind || "unknown";
            if (!groupedData[kind]) {
                groupedData[kind] = {
                    kind: kind,
                    count: 0,
                    items: [],
                    namespace: item.Namespace || "unknown"
                };
            }
            groupedData[kind].count++;
            groupedData[kind].items.push(item);
        });

        const data = Object.values(groupedData).map((group) => {
            const textLength = group.kind.length + group.count.toString().length + 2;
            const baseRadius = (textLength * 6) / 2 + 15;
            const countFactor = Math.sqrt(group.count);
            const estimatedRadius = baseRadius * (1 + (countFactor * 0.15));
            
            return {
                name: toTitleCase(group.kind),
                count: group.count,
                namespace: group.namespace,
                radius: estimatedRadius,
                details: group.items,
                x: Math.random() * window.innerWidth,
                y: Math.random() * window.innerHeight,
            };
        });

        function toTitleCase(str) {
            return str.replace(/\w\S*/g, function(txt) {
                return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
            });
        }

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
            .force("collision", d3.forceCollide().radius(d => d.radius))
            .force("x", d3.forceX(width / 2).strength(0.2))
            .force("y", d3.forceY(height / 2).strength(0.2))
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
                    .attr("r", d.radius * 1.2);
    
                simulation
                    .force("collision", d3.forceCollide().radius(d => 
                        d === d3.select(this).datum() ? d.radius * 1.2 : d.radius
                    ))
                    .alpha(0.5)
                    .restart();
            })
            .on("mouseout", function (event, d) {
                d3.select(this)
                    .transition()
                    .duration(200)
                    .attr("r", d.radius);
    
                simulation
                    .force("collision", d3.forceCollide().radius(d => d.radius))
                    .alpha(0.5)
                    .restart();
            })
            .on("click", function (event, d) {
                zoomIntoBubble(d);
            });
    
        node.append("text")
            .text((d) => `${d.name}: ${d.count}`)
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
                .style("clip-path", `circle(120% at ${bubbleCenterX}px ${bubbleCenterY}px)`)
                .style("-webkit-clip-path", `circle(120% at ${bubbleCenterX}px ${bubbleCenterY}px)`);
        
            const headerContainer = overlay.append("div")
                .attr("class", "header-container")
                .style("display", "flex")
                .style("align-items", "center")
                .style("gap", "10px")
                .style("padding-left", "10px");
        
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
                .style("margin", "0")
                .style("font-family", "Courier, monospace")
                .style("font-size", "24px");
        
            const contentContainer = overlay.append("div")
                .attr("class", "bubble-content")
                .style("padding", "10px");
        
            const itemsList = d.details.map(item => jsyaml.dump(item)).join('\n---\n');
            
            contentContainer.append("pre")
                .style("text-align", "left")
                .style("white-space", "pre-wrap")
                .style("word-break", "break-word")
                .style("font-family", "Courier, monospace")
                .text(itemsList);
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
    updateGraphSize();
});
