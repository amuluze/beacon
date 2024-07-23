function t(t){const e=["B","KB","MB","GB","TB"];let n=0;for(;t>=1024&&n<e.length-1;)t/=1024,n++;return`${t.toFixed(2)} ${e[n]}`}export{t as c};
