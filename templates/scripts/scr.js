async function squat() {
              
    const response = await fetch('http://127.0.0.1:8080/squat',{
        method: 'GET',


    }

    )
    var json = await response.json()
    var xValues = json["x"];
    var yValues = json["y"];
    var barColors = json["color"];
    

    const u = document.getElementById("myChart")
    u.remove()
    const c = document.createElement("canvas")
    c.style="width:100%;max-width:700px;height:500px;"
    c.id = "myChart"
    const g = document.getElementById("graf")
    g.append(c)
    document.getElementById('myChart').innerHTML = new Chart("myChart", {
    type: "bar",
    data: {
      labels: xValues,  
      datasets: [{
        label:"Присед",
        backgroundColor: barColors,
        data: yValues,
        borderWidth: 1
      }]
      
    },
    

});
    console.log(json)
}

async function bench() {
              
    const response = await fetch('http://127.0.0.1:8080/bench',{
        method: 'GET',


    }

    )
      // // добавляем элемент h1 на страницу в элемент body
    var json = await response.json()
    var xValues = json["x"];
    var yValues = json["y"];
     var barColors = json["color"];
     const u = document.getElementById("myChart")
    u.remove()
    const c = document.createElement("canvas")
    c.style="width:100%;max-width:700px;height:500px;"
    c.id = "myChart"
    const g = document.getElementById("graf")
    g.append(c)
    document.getElementById('myChart').innerHTML = new Chart("myChart", {
    type: "bar",
    data: {
      labels: xValues,  
      datasets: [{
        label:"Жим",
        backgroundColor: barColors,
        data: yValues
      }]
    },
    
    

});
    console.log(json)
}

