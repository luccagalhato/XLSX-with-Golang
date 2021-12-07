function criarObjeto() {
  let objeto = {};
  return objeto
}
function getElement(value) {
  return document.getElementById(value)
}
function criarInit(body,method) {
  
  var myInit = criarObjeto();
  if (body == null) {
    myInit.method = method
    myInit.headers = new Headers()
    myInit.mode = 'cors'
    myInit.cache = 'default'
    return myInit
  }
    myInit.method = method
    myInit.headers = new Headers()
    myInit.mode = 'cors'
    myInit.cache = 'default'
    myInit.body = body
    return myInit
}
async function getIds(body, path) {
  var myInit = criarInit(body,'POST')

  const response = await fetch(
    `${window.location.protocol}//${window.location.host}/${path}`,
    myInit
  );
  const blob = await response.blob();
  const text = await blob.text();
  const ids = JSON.parse(text);
  if (response.status != 200) {
    var load = getElement("load");
    var btnenviar = getElement("btnenviar");
    btnenviar.setAttribute(
      "class",
      "row justify-content-center align-items-center"
    );
    load.setAttribute(
      "class",
      "row justify-content-center align-items-center d-none"
    );
    alert("Falha no Download");

    return;
  }
  console.log(ids)
  for (const i in ids) {
    await downloadFile(ids[i]);
  }
}

async function downloadFile(id) {
  var myInit = criarInit(null,'GET')
  console.log(myInit)
  var response = await fetch(
    `${window.location.protocol}//${window.location.host}/id/${id}`,
    myInit
  );
  var myBlob = await response.blob();
  var a = document.createElement("a");
  var url = window.URL.createObjectURL(myBlob);
  a.href = url;
  var filename = response.headers.get("File-Name");
  a.download = filename || "data.xlsx";
  this.console.log(filename);
  a.click();
  a.remove();
  var load = getElement("load");
  var btnenviar = getElement("btnenviar");
  btnenviar.setAttribute(
    "class",
    "row justify-content-center align-items-center"
  );
  load.setAttribute(
    "class",
    "row justify-content-center align-items-center d-none"
  );
  alert("Download Concluído");
  window.URL.revokeObjectURL(url);
}

function enviarData() {
  var input = getElement("dateInicial");
  var input2 = getElement("dateFinal");
  var date = criarObjeto();
  date.dataInicial = input.value
  date.dataFinal = input2.value
  var load = getElement("load");
  var btnenviar = document.getElementById("btnenviar");
  btnenviar.setAttribute(
    "class",
    "row justify-content-center align-items-center d-none"
  );
  load.setAttribute("class", "row justify-content-center align-items-center");

  var datas = JSON.stringify(date);
  getIds(datas, "date");
}

function enviarGtin() {
  var input = getElement("text");
  var load = getElement("load");
  var btnenviar = getElement("btnenviar");
  btnenviar.setAttribute(
    "class",
    "row justify-content-center align-items-center d-none"
  );
  load.setAttribute("class", "row justify-content-center align-items-center");

  getIds(input.value, "gtin");
}

function enviarXML() {
  var input = document.getElementById("file");
  if (input.files.length > 0) {
    let formData = new FormData();
    formData.append("file", input.files[0], input.files[0].name);
    var load = getElement("load");
    var btnenviar = getElement("btnenviar");
    btnenviar.setAttribute(
      "class",
      "row justify-content-center align-items-center d-none"
    );
    load.setAttribute("class", "row justify-content-center align-items-center");

    getIds(formData, "xml");
  }
}
