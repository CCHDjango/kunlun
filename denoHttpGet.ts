const json = fetch("https://stock.xueqiu.com/v5/stock/realtime/quotec.json?symbol=SH600519,SH603288,SH600887,SH600872,SZ000858,SH600600,SZ002302&_=0");

json.then((response) => {
  return response.json();
}).then((jsonData) => {
  console.log(jsonData);
});
