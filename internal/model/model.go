package model

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

const (
	AverageScheduleTime float64 = 5 * 60 * 1000 // 5 minutes in milliseconds
)

type Object struct {
	id            int
	schedule      []float64
	spreadPercent float64
	rs            RandomSupport
}

func NewObject(id int, initialSchedule float64, spreadPercent float64) *Object {
	return &Object{
		id:            id,
		schedule:      []float64{initialSchedule},
		spreadPercent: spreadPercent,
	}
}

func (o *Object) SetRandomSupport(rs RandomSupport) *Object {
	o.rs = rs
	return o
}

func (o *Object) addSchedule(millis float64) *Object {
	o.schedule = append(o.schedule, millis)
	return o
}

func (o *Object) AddRandomSchedule() {
	if len(o.schedule) == 0 {
		panic("No schedules defined")
	}
	lastIdx := len(o.schedule) - 1

	nextSchedule := o.schedule[lastIdx] + AverageScheduleTime
	if o.rs.RandomlyDecide(o.spreadPercent) {
		nextSchedule = o.rs.RandomlyChange(nextSchedule, o.spreadPercent)
	}
	o.addSchedule(nextSchedule)
}

func (o *Object) LastSchedule() float64 {
	if len(o.schedule) == 0 {
		panic("No schedules defined")
	}
	return o.schedule[len(o.schedule)-1]
}

func (o *Object) Schedules() []float64 {
	//return a copy of the schedules
	return append([]float64(nil), o.schedule...)
}

// AsCSVString returns the object representation as a single line in CSV format.
// The first value is the object ID, followed by the consecutive schedule values.
func (o *Object) asCSVString() string {
	bld := strings.Builder{}

	bld.WriteString(strconv.Itoa(o.id))
	for _, schedule := range o.schedule {
		bld.WriteRune(',')
		bld.WriteString(strconv.FormatFloat(schedule, 'f', -1, 64))
	}

	return bld.String()
}

func fromCSVString(line string) (*Object, error) {
	//parse line
	vals := strings.Split(line, ",")
	id, err := strconv.Atoi(vals[0])
	if err != nil {
		return nil, err
	}

	obj := Object{
		id: id,
	}

	for i := 1; i < len(vals); i++ {
		schedule, err := strconv.ParseFloat(vals[i], 64)
		if err != nil {
			return nil, err
		}
		obj.addSchedule(schedule)
	}
	return &obj, nil
}

type ObjSet []*Object

func (oset ObjSet) Marshal(file io.Writer) error {

	bld := strings.Builder{}
	for _, obj := range oset {
		//write CSV line and a newline as a single write operation
		bld.WriteString(obj.asCSVString())
		bld.WriteRune('\n')

		_, err := file.Write([]byte(bld.String()))
		if err != nil {
			return err
		}

		bld.Reset()
	}
	return nil
}

func UnmarshalObjSet(file io.Reader) (ObjSet, error) {

	var res ObjSet = make([]*Object, 0)

	bf := bufio.NewScanner(file)

	for bf.Scan() {
		line := bf.Text()
		obj, err := fromCSVString(line)
		if err != nil {
			return nil, err
		}
		res = append(res, obj)
	}

	return res, nil
}
