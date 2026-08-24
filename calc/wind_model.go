package calc

import "math"

type WindModel struct {
	Speed     float64
	Amplitude float64
}

func NewWindModel(speed, amplitude float64) WindModel {
	if speed < 0 {
		speed = 0
	}
	if amplitude < 0 {
		amplitude = 0
	}
	return WindModel{Speed: speed, Amplitude: amplitude}
}

func (m WindModel) Frequency() float64 {
	frequency := 0.2 + m.Speed/100.0
	if frequency > 8 {
		return 8
	}
	return frequency
}

func (m WindModel) Lateral(t float64) float64 {
	return m.Amplitude * math.Sin(t*m.Frequency())
}

func (m WindModel) Vertical(t float64) float64 {
	return m.Amplitude * 0.35 * math.Cos(t*m.Frequency()*0.7)
}

func (m WindModel) Torsion(t float64) float64 {
	return m.Amplitude * 0.1 * math.Sin(t*m.Frequency()*1.3)
}

func (m WindModel) Energy(duration float64) float64 {
	if duration <= 0 {
		return 0
	}
	return duration * (m.Speed*m.Speed + m.Amplitude*m.Amplitude) / 2
}
