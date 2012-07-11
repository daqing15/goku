package form

import (
    //"math"
    "regexp"
)

type Form struct {
    Name    string
    Action  string
    Fields  map[string]Field
    IsValid bool
    //validResults map[string]*ValidResult
    //hasValid     bool
}

// check all the values on the field are validated,
// if validated, set the clean value
func (fm *Form) Valid() bool {
    isValid := true
    for _, field := range fm.Fields {
        var vr *ValidResult = field.Valid()
        if !vr.IsValid {
            isValid = false
        } else {
            field.CleanValue()
        }
        //fm.validResults[name] = vr
    }
    //fm.hasValid = true
    fm.IsValid = isValid
    return isValid
}

// get the source values
// the values are the same as the submit on the request
func (fm *Form) Values() map[string]string {
    r := make(map[string]string)
    for name, field := range fm.Fields {
        r[name] = field.Value()
    }
    return r
}

// get the clean value after validated
// before you get fm.ClentValues(),
// you must call fm.Valid() first.
func (fm *Form) CleanValues() map[string]interface{} {
    r := make(map[string]interface{})
    for name, field := range fm.Fields {
        r[name] = field.CleanValue()
    }
    return r
}

func (fm *Form) Errors() map[string]string {
    r := make(map[string]string)
    for name, field := range fm.Fields {
        if !field.IsValid() {
            r[name] = field.NickName() + ": " + field.ErrorMsg()
        }
    }
    return r
}

func (fm *Form) FillByMap(m map[string]string) {
    if m == nil {
        return
    }
    for k, v := range m {
        if field, ok := fm.Fields[k]; ok {
            field.SetValue(v)
        }
    }
}

// create a new form whit fields
func NewForm(fields ...Field) *Form {
    f := &Form{
        Fields: make(map[string]Field),
        //validResults: make(map[string]*ValidResult),
    }
    for _, field := range fields {
        f.Fields[field.Name()] = field
    }
    return f
}

type Field interface {
    Name() string
    NickName() string
    ErrorMsg() string
    Value() string
    SetValue(s string)
    CleanValue() interface{}
    IsValid() bool
    Valid() *ValidResult
}

type FieldOption struct {
    Required bool
    NotTrim  bool // not trim the Whitespace
    Range    [2]int
    ErrorMsg string
}

// type Fieldd struct {
//     Name      string
//     Nickname  string
//     Source    string
//     Validater Validater
//     Option    *FieldOption
//     val       interface{}
//     //Val() interface{}
// }

type BaseField struct {
    name       string
    nickname   string
    errorMsg   string
    value      string
    cleanValue interface{}
    isValid    bool
    validater  Validater
    option     *FieldOption
}

func (bf *BaseField) init(name string, nickname string, required bool) {
    bf.name = name
    bf.nickname = nickname
    bf.option = &FieldOption{Required: required}
}

func (bf *BaseField) Field() Field {
    return bf
}

func (bf *BaseField) Name() string {
    return bf.name
}

func (bf *BaseField) NickName() string {
    return bf.nickname
}

func (bf *BaseField) ErrorMsg() string {
    return bf.errorMsg
}

func (bf *BaseField) Value() string {
    return bf.value
}

func (bf *BaseField) SetValue(val string) {
    bf.value = val
}

func (bf *BaseField) CleanValue() interface{} {
    return bf.cleanValue
}

func (bf *BaseField) IsValid() bool {
    return bf.isValid
}

func (bf *BaseField) Valid() *ValidResult {
    vr := bf.validater.Valid(bf.value, bf.option)
    bf.isValid = vr.IsValid
    bf.cleanValue = vr.CleanValue
    bf.errorMsg = vr.ErrorMsg
    return vr
}

func (bf *BaseField) Required(r bool) *BaseField {
    bf.option.Required = r
    return bf
}

func (bf *BaseField) Max(max int) *BaseField {
    bf.option.Range[1] = max
    return bf
}

func (bf *BaseField) Min(min int) *BaseField {
    bf.option.Range[0] = min
    return bf
}

func (bf *BaseField) MaxLength(max int) *BaseField {
    bf.option.Range[1] = max
    return bf
}

func (bf *BaseField) MinLength(min int) *BaseField {
    bf.option.Range[0] = min
    return bf
}

func (bf *BaseField) Range(min int, max int) *BaseField {
    bf.option.Range = [2]int{min, max}
    return bf
}

type IntegerField struct {
    BaseField
}

type CharField struct {
    BaseField
}

type TextField struct {
    BaseField
}

type RegexpField struct {
    BaseField
    regexp *regexp.Regexp
}

func NewCharField(name string, nickname string, required bool) *CharField {
    tf := &CharField{}
    tf.init(name, nickname, required)
    tf.validater = &stringValidater{}
    return tf
}

func NewTextField(name string, nickname string, required bool) *TextField {
    tf := &TextField{}
    tf.init(name, nickname, required)
    tf.validater = &stringValidater{}
    return tf
}

func NewRegexpField(name string, nickname string, required bool, re string) *RegexpField {
    tf := &RegexpField{}
    tf.init(name, nickname, required)
    tf.regexp = regexp.MustCompile(re)
    tf.validater = &regexpValidater{
        Regexp: tf.regexp,
    }
    return tf
}

func NewIntegerField(name string, nickname string, required bool) *IntegerField {
    tf := &IntegerField{}
    tf.init(name, nickname, required)
    tf.validater = &intValidater{}
    return tf
}
