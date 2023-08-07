package accounts

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/domain/object"
)

// Request body for `POST /v1/accounts`
type AddRequest struct {
	Username string
	Password string
}

// Handle request for `POST /v1/accounts`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	/*
		ハンドラが抽象化を介してドメインを介さず直接DBを触っていたので、一部をドメインでやるにしても、
		最も上流のロジックは全てハンドラに入らざるを得ないので、
		ハンドラの役目であるリクエストのバリデーションや認証以外の仕組みがごちゃ混ぜになってしまうように感じた
		本質的にMVCにおけるFat Controller問題と同様に「ネットワークの仕様に沿った入力をドメインに変換する処理」と
		「ドメインの本質的な処理」が混ざり見通しが悪くなる問題を感じた。
		実際、今回の実装ではdomain.objectパッケージの機能がかなり貧弱になってしまった

		ハンドラに最上位の処理をやらせないという点で言うと、ドメインを独立させるのがベターな実装かもしれない
		現在はハンドラがDAOをインターフェースを介して呼び出して、domain.objectを生成して
		そのメソッドを呼び出すことでなんらかの処理を行い、DAOに保存をさせると言う構成で、処理の主体がハンドラである。
		これを修正して、まずハンドラは入力のパース検証メソッド呼び分けに専念させて、データが揃ったら
		ドメインの入り口に対して整形済みデータを引数にしてコマンドを呼び出して、そもそもDBへのインターフェースを使わない
		ドメインは与えられた綺麗なデータをもとにdomain.objectがdomain.repositoryを介してdaoを呼び出しつつ処理を行う方式にする
		domain.repositoryのインターフェースによりドメインはDBに依存しない
		（多分domain.repositoryにインターフェースがあるのはドメインから使って欲しいという気持ちからだと思う）
		（現在の構成だと、ドメインにインターフェースが置いてあるだけで、実装も使用も全くdomain.objectでは行われていないので、単にdomain配下に置いただけになっている）
		また、domain.objectのメソッドをハンドラが呼び出すが、その逆はないため、ドメインはハンドラにも依存しない
	*/
	ctx := r.Context()

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account := new(object.Account)
	account.Username = req.Username
	if err := account.SetPassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	account.SetCreateAt()

	account, err := h.ar.CreateNewAccount(ctx, account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
