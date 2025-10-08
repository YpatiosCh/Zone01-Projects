function getAcceleration(obj) {
  if (typeof obj.f === 'number' && typeof obj.m === 'number' && obj.m !== 0) {
    return obj.f / obj.m;
  }
  if (typeof obj["Δv"] === 'number' && typeof obj["Δt"] === 'number' && obj["Δt"] !== 0) {
    return obj["Δv"] / obj["Δt"];
  }
  if (typeof obj.d === 'number' && typeof obj.t === 'number' && obj.t !== 0) {
    return (2 * obj.d) / (obj.t ** 2);
  }
  return "impossible";
}
