package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-numb/go-ftx/auth"
	"github.com/go-numb/go-ftx/rest"
	"github.com/go-numb/go-ftx/rest/private/spotmargin"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  "run sample command",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		ftxKey := viper.GetString("ftx.key")
		ftxSecret := viper.GetString("ftx.secret")
		if ftxKey == "xxx" || ftxSecret == "xxx" {
			log.Fatal("please setup your config.yml file")
		}
		crontab := cron.New()

		// client := rest.New(auth.New(ftxKey, ftxSecret))
		// or
		// UseSubAccounts
		clientWithSubAccounts := rest.New(
			auth.New(
				ftxKey,
				ftxSecret,
				// auth.SubAccount{
				// 	UUID:     1,
				// 	Nickname: "gunbot",
				// },
			))
		// clientWithSubAccounts.Auth.UseSubAccountID(1) // or 2... this number is key in map[int]SubAccount

		// client or clientWithSubAccounts in this time.
		c := clientWithSubAccounts // or clientWithSubAccounts

		coin := "USD"
		updateLending(c, coin)
		_, err := crontab.AddFunc("5 * * * *", func() { updateLending(c, coin) })
		if err != nil {
			log.Fatalf("starting cron: %s", err)
		}
		crontab.Start()

		<-ctx.Done()
		crontab.Stop()
		log.Info("got stop signal, exiting")
	},
}

func updateLending(c *rest.Client, coin string) {

	rli := spotmargin.RequestForLendingInfo{}
	lending, err := c.GetLendingInfo(&rli)
	if err != nil || lending == nil {
		log.Fatalf("loading lending info: %s", err)
	}

	for _, li := range *lending {
		if li.Coin == coin {
			moolah := li.Lendable - li.Locked
			if moolah > 0.01 {
				rlo := spotmargin.RequestForLendingOffer{
					Coin: li.Coin,
					Size: float64(int(li.Lendable*100)) / 100,
					Rate: li.MinRate,
				}
				_, err := c.SubmitLendingOffer(&rlo)
				if err != nil {
					spew.Dump(li, rlo)
					log.Fatalf("failed to submit lending offer: %s", err)
				}
				apr := moolah * 24 * 365 / li.Lendable // maybe ok if ran every hour
				log.WithField("APR", apr).Infof("added $%.2f to lending for a grand total of $%.2f", moolah, li.Lendable)
			} else {
				log.Info("nothing to increase offer with")
			}
		}
	}
}

func init() {
	RootCmd.AddCommand(runCmd)

	addRunFlags(runCmd)
}
