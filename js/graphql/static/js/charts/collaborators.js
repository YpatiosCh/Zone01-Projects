const SVG_NS = "http://www.w3.org/2000/svg";

export function createCollaboratorChart(data) {
  if (!Array.isArray(data) || data.length === 0) {
    return null;
  }

  const section = document.createElement("section");
  section.className = "collab-chart";

  const title = document.createElement("h3");
  title.className = "collab-chart__title";
  title.textContent = "Collaborators";
  section.appendChild(title);

  const description = document.createElement("p");
  description.className = "collab-chart__subtitle";
  description.textContent = "Number of shared group projects";
  section.appendChild(description);

  const svg = buildChartSvg(data);
  section.appendChild(svg);

  return section;
}

function buildChartSvg(data) {
  const width = 700;
  const barHeight = 32;
  const barGap = 16;
  const margin = { top: 32, right: 36, bottom: 36, left: 160 };
  const innerHeight = data.length * (barHeight + barGap) - barGap;
  const height = margin.top + Math.max(innerHeight, barHeight) + margin.bottom;
  const chartWidth = width - margin.left - margin.right;

  const svg = document.createElementNS(SVG_NS, "svg");
  svg.setAttribute("viewBox", `0 0 ${width} ${height}`);
  svg.setAttribute("preserveAspectRatio", "xMidYMid meet");
  svg.setAttribute("role", "img");
  svg.setAttribute(
    "aria-label",
    "Bar chart showing collaborator frequency by username",
  );

  const maxCount = Math.max(...data.map((d) => d.count));
  const scale = maxCount > 0 ? chartWidth / maxCount : 0;

  const background = document.createElementNS(SVG_NS, "rect");
  background.setAttribute("x", margin.left);
  background.setAttribute("y", margin.top);
  background.setAttribute("width", chartWidth);
  background.setAttribute("height", Math.max(innerHeight, barHeight));
  background.setAttribute("fill", "rgba(142, 208, 255, 0.08)");
  background.setAttribute("rx", "12");
  svg.appendChild(background);

  data.forEach((entry, index) => {
    const barGroup = document.createElementNS(SVG_NS, "g");

    const y = margin.top + index * (barHeight + barGap);
    const barWidth = scale * entry.count;

    const label = document.createElementNS(SVG_NS, "text");
    label.setAttribute("x", margin.left - 16);
    label.setAttribute("y", y + barHeight / 2);
    label.setAttribute("text-anchor", "end");
    label.setAttribute("dominant-baseline", "middle");
    label.setAttribute("class", "collab-chart__label");
    label.textContent = entry.login;
    barGroup.appendChild(label);

    const rect = document.createElementNS(SVG_NS, "rect");
    rect.setAttribute("x", margin.left);
    rect.setAttribute("y", y);
    rect.setAttribute("width", Math.max(barWidth, 8));
    rect.setAttribute("height", barHeight);
    rect.setAttribute("rx", "10");
    rect.setAttribute("class", "collab-chart__bar");
    rect.setAttribute("aria-hidden", "true");
    barGroup.appendChild(rect);

    const value = document.createElementNS(SVG_NS, "text");
    const textOffset = 12;
    const labelInside = barWidth > chartWidth * 0.66;
    const textX = labelInside
      ? margin.left + Math.max(barWidth - 8, 12)
      : Math.min(
          margin.left + barWidth + textOffset,
          margin.left + chartWidth - 4,
        );
    value.setAttribute("x", textX);
    value.setAttribute("y", y + barHeight / 2);
    value.setAttribute("class", "collab-chart__value");
    value.setAttribute("dominant-baseline", "middle");
    value.setAttribute("text-anchor", labelInside ? "end" : "start");
    if (labelInside) {
      value.classList.add("collab-chart__value--inside");
    }
    value.textContent = `${entry.count}`;
    barGroup.appendChild(value);

    svg.appendChild(barGroup);
  });

  const axis = document.createElementNS(SVG_NS, "line");
  axis.setAttribute("x1", margin.left);
  axis.setAttribute("y1", margin.top + Math.max(innerHeight, barHeight));
  axis.setAttribute("x2", margin.left + chartWidth);
  axis.setAttribute("y2", margin.top + Math.max(innerHeight, barHeight));
  axis.setAttribute("class", "collab-chart__axis");
  svg.appendChild(axis);

  return svg;
}
