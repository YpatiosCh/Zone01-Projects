import { formatXPAmount } from "../utils/profile.js";

const SVG_NS = "http://www.w3.org/2000/svg";

export function createProgressChart(points) {
  if (!Array.isArray(points) || points.length === 0) {
    return null;
  }

  const section = document.createElement("section");
  section.className = "progress-chart";

  const title = document.createElement("h3");
  title.className = "progress-chart__title";
  title.textContent = "XP Progress Over Time";
  section.appendChild(title);

  const subtitle = document.createElement("p");
  subtitle.className = "progress-chart__subtitle";
  subtitle.textContent = "Cumulative XP earned per date";
  section.appendChild(subtitle);

  const svg = buildProgressSvg(points);
  section.appendChild(svg);

  return section;
}

function buildProgressSvg(points) {
  const width = 720;
  const height = 380;
  const margin = { top: 40, right: 40, bottom: 60, left: 60 };
  const plotWidth = width - margin.left - margin.right;
  const plotHeight = height - margin.top - margin.bottom;

  const svg = document.createElementNS(SVG_NS, "svg");
  svg.setAttribute("viewBox", `0 0 ${width} ${height}`);
  svg.setAttribute("preserveAspectRatio", "xMidYMid meet");
  svg.setAttribute("role", "img");
  svg.setAttribute(
    "aria-label",
    "Line chart showing cumulative XP progress over time",
  );

  const cumulativePoints = computeCumulative(points);
  if (cumulativePoints.length === 0) {
    return svg;
  }

  const xScale = createTimeScale(
    cumulativePoints.map((d) => d.date),
    plotWidth,
  );
  const yScale = createLinearScale(
    cumulativePoints.map((d) => d.value),
    plotHeight,
  );

  const chartGroup = document.createElementNS(SVG_NS, "g");
  chartGroup.setAttribute("transform", `translate(${margin.left},${margin.top})`);
  svg.appendChild(chartGroup);

  const axes = drawAxes(chartGroup, plotWidth, plotHeight, xScale, yScale);
  axes.forEach((el) => chartGroup.appendChild(el));

  const gridLines = drawGrid(plotWidth, plotHeight, yScale.ticks);
  gridLines.forEach((line) => chartGroup.appendChild(line));

  const path = document.createElementNS(SVG_NS, "path");
  path.setAttribute("class", "progress-chart__line");
  path.setAttribute("d", buildPath(cumulativePoints, xScale, yScale, plotHeight));
  chartGroup.appendChild(path);

  cumulativePoints.forEach((point) => {
    const cx = xScale.map(point.date);
    const cy = plotHeight - yScale.map(point.value);

    const dot = document.createElementNS(SVG_NS, "circle");
    dot.setAttribute("class", "progress-chart__dot");
    dot.setAttribute("cx", cx);
    dot.setAttribute("cy", cy);
    dot.setAttribute("r", 4);
    chartGroup.appendChild(dot);
  });

  return svg;
}

function computeCumulative(points) {
  let accumulator = 0;
  return points
    .filter((point) => point && point.date instanceof Date && point.value >= 0)
    .sort((a, b) => a.date - b.date)
    .map((point) => {
      accumulator += point.value;
      return {
        date: point.date,
        value: accumulator,
      };
    });
}

function createTimeScale(dates, width) {
  const minDate = dates[0];
  const maxDate = dates[dates.length - 1];
  const range = maxDate.getTime() - minDate.getTime() || 1;

  const scale = (date) =>
    ((date.getTime() - minDate.getTime()) / range) * width;

  return {
    map: (date) => scale(date),
    ticks: generateDateTicks(minDate, maxDate, width),
  };
}

function generateDateTicks(minDate, maxDate, width) {
  const desiredTickCount = Math.max(3, Math.floor(width / 120));
  const ticks = [];

  const totalMs = maxDate.getTime() - minDate.getTime();
  if (totalMs === 0) {
    return [{ date: minDate, label: formatDateLabel(minDate), position: 0 }];
  }

  const step = totalMs / (desiredTickCount - 1);
  for (let i = 0; i < desiredTickCount; i += 1) {
    const tickDate = new Date(minDate.getTime() + step * i);
    ticks.push({
      date: tickDate,
      label: formatDateLabel(tickDate),
      position: (step * i) / totalMs,
    });
  }

  return ticks;
}

function createLinearScale(values, height) {
  const maxValue = Math.max(...values, 0);
  const range = maxValue || 1;

  return {
    map: (value) => (value / range) * height,
    range,
    ticks: generateValueTicks(maxValue),
  };
}

function generateValueTicks(maxValue) {
  if (maxValue <= 0) {
    return [];
  }
  const tickCount = 5;
  const step = maxValue / tickCount;
  const ticks = [];
  for (let i = 0; i <= tickCount; i += 1) {
    const value = step * i;
    ticks.push({
      value,
      label: formatXPAmount(value),
      position: i / tickCount,
    });
  }
  return ticks;
}

function formatDateLabel(date) {
  return date.toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
  });
}

function buildPath(points, xScale, yScale, plotHeight) {
  return points
    .map((point, index) => {
      const x = xScale.map(point.date);
      const y = plotHeight - yScale.map(point.value);
      return `${index === 0 ? "M" : "L"}${x},${y}`;
    })
    .join(" ");
}

function drawAxes(group, width, height, xScale, yScale) {
  const elements = [];

  // X axis
  const xAxis = document.createElementNS(SVG_NS, "line");
  xAxis.setAttribute("class", "progress-chart__axis");
  xAxis.setAttribute("x1", 0);
  xAxis.setAttribute("y1", height);
  xAxis.setAttribute("x2", width);
  xAxis.setAttribute("y2", height);
  elements.push(xAxis);

  xScale.ticks.forEach((tick) => {
    const x = tick.position * width;
    const tickLine = document.createElementNS(SVG_NS, "line");
    tickLine.setAttribute("class", "progress-chart__tick");
    tickLine.setAttribute("x1", x);
    tickLine.setAttribute("y1", height);
    tickLine.setAttribute("x2", x);
    tickLine.setAttribute("y2", height + 8);
    elements.push(tickLine);

    const label = document.createElementNS(SVG_NS, "text");
    label.setAttribute("class", "progress-chart__axis-label");
    label.setAttribute("x", x);
    label.setAttribute("y", height + 24);
    label.setAttribute("text-anchor", "middle");
    label.textContent = tick.label;
    elements.push(label);
  });

  // Y axis
  const yAxis = document.createElementNS(SVG_NS, "line");
  yAxis.setAttribute("class", "progress-chart__axis");
  yAxis.setAttribute("x1", 0);
  yAxis.setAttribute("y1", 0);
  yAxis.setAttribute("x2", 0);
  yAxis.setAttribute("y2", height);
  elements.push(yAxis);

  yScale.ticks.forEach((tick) => {
    const y = height - tick.position * height;
    const tickLine = document.createElementNS(SVG_NS, "line");
    tickLine.setAttribute("class", "progress-chart__tick");
    tickLine.setAttribute("x1", -8);
    tickLine.setAttribute("y1", y);
    tickLine.setAttribute("x2", 0);
    tickLine.setAttribute("y2", y);
    elements.push(tickLine);

    const label = document.createElementNS(SVG_NS, "text");
    label.setAttribute("class", "progress-chart__axis-label");
    label.setAttribute("x", -12);
    label.setAttribute("y", y + 4);
    label.setAttribute("text-anchor", "end");
    label.textContent = tick.label;
    elements.push(label);
  });

  return elements;
}

function drawGrid(width, height, ticks) {
  return ticks.map((tick) => {
    const y = height - tick.position * height;
    const line = document.createElementNS(SVG_NS, "line");
    line.setAttribute("class", "progress-chart__grid-line");
    line.setAttribute("x1", 0);
    line.setAttribute("y1", y);
    line.setAttribute("x2", width);
    line.setAttribute("y2", y);
    return line;
  });
}
