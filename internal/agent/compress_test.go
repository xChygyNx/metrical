package agent

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	messages := [][]byte{
		[]byte("Hello, how are you?"),
		[]byte("skhg;ah;lk;kha;ovjap'dhfia o vipi oypo poapih opa oig og ihaoiuvg ppayg pi p\n" +
			"321 61461 64+ 11 64 64 +94 646 894 6849498464+9481684646 6464169846168 1916463\n" +
			"*^&^*(&$(%^*&^&^%^(&^^*YU()(_()&*%*^##&*^(&^%*(^*^$&^&&&&^&^*&*&^&$^&%&^&%$#&^"),
	}
	for _, message := range messages {
		t.Run(fmt.Sprintf("compress %s", message[:15]), func(t *testing.T) {
			res, err := compress(message)
			assert.Nil(t, err)
			assert.NotEqual(t, len(res), len(message), fmt.Sprintf("%d !< %d", len(res), len(message)))
		})
	}
}
