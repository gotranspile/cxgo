package cxgo

import "testing"

var casesRun = []struct {
	name string
	src  string
}{
	{
		name: "issue 77",
		src: `
#include <stdio.h>

int globe = 35;


int* fn1(int strd) {
	static int result = 0; 
    result = globe;
    if (strd > 2) {
        result += 4;
    }
    result += strd;
    return &result; 
}

int main() {
    char q_q = 34;
    int p_p = 110; 
    p_p += q_q ^ p_p || globe; 
    q_q = globe ^ p_p ^ 25;
    printf("%d\n", *fn1(q_q));
    printf("%d\n", *fn1(p_p));
	int t1 = *fn1(q_q);
	int t2 = *fn1(p_p);
    int final=t1 += t2; 
    printf("Result: %d\n", final ); 
    return 0;
}
`,
	},
}

func TestRunTranslated(t *testing.T) {
	for _, c := range casesRun {
		c := c
		t.Run(c.name, func(t *testing.T) {
			testTranspileOut(t, c.src)
		})
	}
}
