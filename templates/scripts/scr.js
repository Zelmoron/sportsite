function starterBar(){
  document.getElementById('myChart').innerHTML = new Chart("myChart", {
    type: "bar",
    data: {
      labels: ['/',"/",'/','/',"/",'/','/',"/",'/','/',"/",'/'],  
      datasets: [{
        label:"Menu",
        backgroundColor: ["gray"],
        data: [1,5,9,14,18,22,26,30,34,38,42,46],
        borderWidth: 1
      }]
      
    },
    
    
  
  });
}
starterBar()

var s = 0 
function getType(value){
  
  s+=value
  

  switch (s){
    case 0:
      return "bar"
    case 1:
      return "line"
    case 2:
      return "doughnut"
    case 3:
      return "pie"
    case 4:
      return "polarArea"
    case 5:
      return "radar"
    case 6:
      s = 0
      return "bar"
      
  }
  
  
}

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
    type: getType(0),
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
    type: getType(0),
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

async function dead() {
              
  const response = await fetch('http://127.0.0.1:8080/dead',{
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
  type: getType(0),
  data: {
    labels: xValues,  
    datasets: [{
      label:"Становая тяга",
      backgroundColor: barColors,
      data: yValues
    }]
  },
  
  

});
  console.log(json)
}


async function pull() {
              
  const response = await fetch('http://127.0.0.1:8080/pull',{
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
  type: getType(0),
  data: {
    labels: xValues,  
    datasets: [{
      label:"Подтягивания",
      backgroundColor: barColors,
      data: yValues
    }]
  },
  
  

});
  console.log(json)
}


async function ton() {
              
  const response = await fetch('http://127.0.0.1:8080/ton',{
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
  type: getType(0),
  data: {
    labels: xValues,  
    datasets: [{
      label:"Тонаж",
      backgroundColor: barColors,
      data: yValues
    }]
  },
  
  

});
  console.log(json)
}


