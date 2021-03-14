package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2"
)

type d3Struct struct {
	Name    string   `json:"name"`
	Parent  string   `json:"parent"`
	Value   *float64 `json:"value"`
	Percent *float64 `json:"percent"`
}

func pFloat(f float64) *float64 {
	return &f
}

func chartBalance(client *binance.Client, useUSDT bool) {
	rawBalances := mustRequestAccountBalances(client)
	ticker := mustRequestTicker(client)
	balances := mustCalculateAllBalances(rawBalances, ticker, false)
	mustPrintChartHTML(balances, useUSDT)
}

func mustPrintChartHTML(balances map[string]assetStatus, useUSDT bool) {
	data := []d3Struct{{Name: "Origin", Parent: "", Value: nil, Percent: nil}}
	for asset, balance := range balances {
		if asset == "Total" {
			continue
		}
		value := pFloat(balance.BalanceInBTC)
		if useUSDT {
			value = pFloat(balance.BalanceInUSDT)
		}
		data = append(data, d3Struct{
			Name:    asset,
			Parent:  "Origin",
			Value:   value,
			Percent: pFloat(balance.Delta24pct),
		})
	}

	jsonMarshalledData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(`
<!DOCTYPE html>
<html>
<head>   
  <meta charset="utf-8">
  <script src="https://d3js.org/d3.v4.min.js"></script>
  <script>
	const data = 
  ` + string(jsonMarshalledData) + `
  </script>
  <style>
	body {
		font-family: 'Lucida Sans', 'Lucida Sans Regular', 'Lucida Grande', 'Lucida Sans Unicode', Geneva, Verdana, sans-serif;
	}
  	.text { 
		  font-size: 16px;
		  fill: white;
		  text-shadow: 1px 1px 2px black;
	}
	</style>
</head>
<body>
	<div id="my_dataviz"></div>
	<ul>
		<li>Values are in BTC.</li>
		<li>Some coins don't have a direct market against BTC. If they can be proxied via BNB, the proxy is used.</li>
		<li>Percentages are the 24-hour delta of each coin against BTC.</li>
		<li>For BTC itself, the % is against USDT.</li>
		<li>Some coins have no direct percentage against BTC. In this case, no % is shown unless it can be proxied.</li>
	</ul>
	<script>
		const chartHeight = 700
		const squarePadding = 1
		const strokeWidth = 1

		const chartWidth = 4/3 * chartHeight

		// set the dimensions and margins of the graph
		var margin = {top: 10, right: 10, bottom: 10, left: 10},
		  width = chartWidth - margin.left - margin.right,
		  height = chartHeight - margin.top - margin.bottom;
		
		// append the svg object to the body of the page
		var svg = d3.select("#my_dataviz")
		.append("svg")
		  .attr("width", width + margin.left + margin.right)
		  .attr("height", height + margin.top + margin.bottom)
		.append("g")
		  .attr("transform",
				"translate(" + margin.left + "," + margin.top + ")");
		
		// stratify the data: reformatting for d3.js
		var root = d3.stratify()
		.id(function(d) { return d.name; })   // Name of the entity (column name is name in csv)
		.parentId(function(d) { return d.parent; })   // Name of the parent (column name is parent in csv)
		(data);
		root.sum(function(d) { return +d.value })   // Compute the numeric value for each entity
	
		// Then d3.treemap computes the position of each element of the hierarchy
		// The coordinates are added to the root object above
		d3.treemap()
		.size([width, height])
		.padding(squarePadding)
		(root)
	
		// use this information to add rectangles:
		svg
		.selectAll("rect")
		.data(root.leaves())
		.enter()
		.append("rect")
			.attr('x', function (d) { return d.x0; })
			.attr('y', function (d) { return d.y0; })
			.attr('width', function (d) { return d.x1 - d.x0; })
			.attr('height', function (d) { return d.y1 - d.y0; })
			.style("stroke", "black")
			.style("stroke-width", ` + "`" + `${strokeWidth}px` + "`" + `)
			.style("fill", function (d) {
				return !d.data.percent ? ` + "`" + `#CCC` + "`" + ` : d.data.percent >= 0
					? ` + "`" + `rgb(151, ${Math.max(230-d.data.percent*4, 0)}, 0)` + "`" + `
					: ` + "`" + `rgb(${Math.max(255-d.data.percent*4, 0)}, 25, 60)` + "`" + `
			});
	
		const urlParams = new URLSearchParams(window.location.search)
		const isPercent = !urlParams.get('p')
		const totalBtc = data.slice(1).map(d => d.value).reduce((a, b) => a + b, 0)

		// and to add the text labels
		svg
		.selectAll("text")
		.data(root.leaves())
		.enter()
		.append("text")
			.attr("x", function(d){ return (d.x1+d.x0)/2})
			.attr("y", function(d){ return (d.y1+d.y0)/2 })
			.attr("dominant-baseline", "middle")
			.attr("text-anchor", "middle")
			.text(function(d){
				const percent = !d.data.percent ? '' : d.data.percent >= 0 
					? ` + "`" + ` (+${d.data.percent.toFixed(0)}%)` + "`" + ` 
					: ` + "`" + ` (${d.data.percent.toFixed(0)}%)` + "`" + `
				return ` + "`" + `${d.data.name} ${isPercent ? ` + "`" + `${Math.round(d.data.value / totalBtc * 100)}%` + "`" + ` : d.data.value.toFixed(2)}${percent}` + "`" + `
			})
			.attr("class", "text")

		</script>
	
</body>
</html>
`)
}
