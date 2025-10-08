function defaultCurry(obj1) {
  return function (obj2) {
    return { ...obj1, ...obj2 };
  };
}

function mapCurry(fn) {
  return function (obj) {
    return Object.fromEntries(Object.entries(obj).map(fn));
  };
}

function reduceCurry(fn) {
  return function (obj, init) {
    return Object.entries(obj).reduce(fn, init);
  };
}

function filterCurry(fn) {
  return function (obj) {
    return Object.fromEntries(Object.entries(obj).filter(fn));
  };
}

function reduceScore(personnel, init) {
  const forceUsers = filterCurry(([_k, v]) => v.isForceUser)(personnel);
  return reduceCurry((acc, [_k, v]) => acc + v.pilotingScore + v.shootingScore)(forceUsers, init);
}


function filterForce(personnel) {
  return filterCurry(([_, v]) => v.isForceUser && v.shootingScore >= 80)(personnel);
}

function mapAverage(personnel) {
  return mapCurry(([k, v]) => {
    const average = (v.pilotingScore + v.shootingScore) / 2;
    return [k, { ...v, averageScore: average }];
  })(personnel);
}
