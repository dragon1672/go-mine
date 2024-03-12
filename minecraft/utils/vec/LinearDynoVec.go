package vec

import "time"

type LinearDynoVec struct {
	snapshotTime time.Time
	lastSnapshot Vec3
	speed        Vec3
}

func (l *LinearDynoVec) Get(t time.Time) Vec3 {
	dt := l.snapshotTime.Sub(t)
	return l.lastSnapshot.Add(l.speed.Mul(dt.Seconds()))
}

func (l *LinearDynoVec) Set(t time.Time, v Vec3, speed Vec3) {
	l.snapshotTime = t
	l.lastSnapshot = v
	l.speed = speed
}

func MakeLinearDynoVec(t time.Time, v Vec3, speed Vec3) *LinearDynoVec {
	l := &LinearDynoVec{}
	l.Set(t, v, speed)
	return l
}
