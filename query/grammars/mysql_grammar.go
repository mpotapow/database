package grammars

type MysqlGrammar struct {
	*Grammar
}

func NewMysqlGrammar() *MysqlGrammar {

	var mg = &MysqlGrammar{
		Grammar: NewGrammar(),
	}

	mg.Grammar.SetParametrizeSymbol("?")
	mg.Grammar.SetSelectComponents(mg.Grammar.GetDefaultSelectComponents())

	return mg
}