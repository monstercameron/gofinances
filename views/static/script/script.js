console.log("script.js loaded");

document.addEventListener("DOMContentLoaded", function () {
  // Function to create SVG element
  function createSvgElement() {
    var svgNS = "http://www.w3.org/2000/svg";
    var svg = document.createElementNS(svgNS, "svg");
    svg.setAttribute("xmlns", svgNS);
    svg.setAttribute("fill", "none");
    svg.setAttribute("viewBox", "0 0 24 24");
    svg.setAttribute("stroke-width", "1.5");
    svg.setAttribute("stroke", "currentColor");
    svg.classList.add("w-8", "h-5", "inline");

    var path = document.createElementNS(svgNS, "path");
    path.setAttribute("stroke-linecap", "round");
    path.setAttribute("stroke-linejoin", "round");
    path.setAttribute("d", "M4.5 10.5 12 3m0 0 7.5 7.5M12 3v18");
    svg.appendChild(path);

    return svg;
  }

  // Select all sortable columns
  var sortableColumns = document.querySelectorAll(".sortable");
  var svgElement = createSvgElement(); // Create the SVG element

  // Function to handle click event on sortable columns
  function handleColumnClick(event) {
    // Remove 'sorted' class from all sortable columns and SVG from its current parent
    sortableColumns.forEach(function (column) {
      column.classList.remove("sorted");
      if (column.contains(svgElement)) {
        column.removeChild(svgElement);
      }
    });

    // Add 'sorted' class to clicked column and append the SVG to it
    event.currentTarget.classList.add("sorted");
    event.currentTarget.appendChild(svgElement);
  }

  // Add click event listener to each sortable column
  sortableColumns.forEach(function (column) {
    column.addEventListener("click", handleColumnClick);
  });
});
