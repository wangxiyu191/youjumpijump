package jump

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

var r *Request

func init() {
	r = NewRequest()
}

func NewTrain(width, height int, ratio float64, phoneOS string) *Train {
	train := &Train{
		width:        width,
		height:       height,
		distances:    []float64{},
		ratios:       map[float64]float64{},
		defaultRatio: ratio,
		phoneOS:      phoneOS,
	}

	log.Printf("Get Train data")
	var body []byte
	var err error
	if phoneOS == "Android" {
		body, err = exec.Command("/system/bin/curl", fmt.Sprintf("http://youjumpijump.faceair.me/%d/%d", width, height)).Output()
	} else if phoneOS == "iOS" {
		_, body, err = r.Get(fmt.Sprintf("http://youjumpijump.faceair.me/%s/%d/%d", phoneOS, width, height))
	}
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(body))
		for scanner.Scan() {
			line := strings.Split(scanner.Text(), ",")
			if len(line) == 2 {
				distance, err1 := strconv.ParseFloat(line[0], 64)
				ratio, err2 := strconv.ParseFloat(line[1], 64)
				if err1 == nil && err2 == nil {
					train.distances = append(train.distances, distance)
					train.ratios[distance] = ratio
				}
			}
		}
	}

	return train
}

type Train struct {
	width        int
	height       int
	distances    []float64
	ratios       map[float64]float64
	defaultRatio float64
	phoneOS      string
}

func (s *Train) Add(distance, ratio float64) {
	if s.phoneOS == "Android" {
		exec.Command("/system/bin/curl", "--data", fmt.Sprintf("%v,%v", distance, ratio), fmt.Sprintf("http://youjumpijump.faceair.me/%s/%d/%d", s.phoneOS, s.width, s.height)).Run()
	} else if s.phoneOS == "iOS" {
		r.Post(fmt.Sprintf("http://youjumpijump.faceair.me/%s/%d/%d", s.phoneOS, s.width, s.height), nil, strings.NewReader(fmt.Sprintf("%v,%v", distance, ratio)))
	}

	s.distances = append(s.distances, distance)
	s.ratios[distance] = ratio
}

func (s *Train) Find(nowDistance float64) (similarDistance, simlarRatio float64) {
	sumR := 0.0
	sumD := 0.0
	count := 0.0

	for _, distance := range s.distances {
		if math.Abs(nowDistance-distance) < 10 {
			count++
			sumD += distance
			sumR += s.ratios[distance]
		}
	}
	if count < 3 {
		return 0, s.defaultRatio
	}

	return sumD / count, sumR / count
}
