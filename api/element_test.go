package api_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/api/internal/mocks"
)

var _ = Describe("Element", func() {
	var (
		element *Element
		session *Session
		bus     *mocks.Bus
		err     error
	)

	BeforeEach(func() {
		bus = &mocks.Bus{}
		session = &Session{bus}
		element = &Element{"some-id", session}
	})

	ItShouldMakeAnElementRequest := func(endpoint, method string, body ...string) {
		It("should hit the desired endpoint", func() {
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/" + endpoint))
		})

		It("should make a "+method+" request", func() {
			Expect(bus.SendCall.Method).To(Equal(method))
		})

		It("should not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		if len(body) > 0 {
			It("should set the request body", func() {
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(body[0]))
			})
		}
	}

	Describe("#Send", func() {
		var result string

		BeforeEach(func() {
			bus.SendCall.Result = `"some result"`
			err = element.Send("endpoint", "method", "body", &result)
		})

		ItShouldMakeAnElementRequest("endpoint", "method", `"body"`)

		It("should retrieve the result", func() {
			Expect(result).To(Equal("some result"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				err := element.Send("endpoint", "method", "body", &result)
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElement", func() {
		var singleElement *Element

		BeforeEach(func() {
			bus.SendCall.Result = `{"ELEMENT": "some-id"}`
			singleElement, err = element.GetElement(Selector{"css selector", "#selector"})
		})

		ItShouldMakeAnElementRequest("element", "POST", `{"using": "css selector", "value": "#selector"}`)

		It("should return an element with an ID and session", func() {
			Expect(singleElement.ID).To(Equal("some-id"))
			Expect(singleElement.Session).To(Equal(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetElement(Selector{"css selector", "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElements", func() {
		var elements []*Element

		BeforeEach(func() {
			bus.SendCall.Result = `[{"ELEMENT": "some-id"}, {"ELEMENT": "some-other-id"}]`
			elements, err = element.GetElements(Selector{"css selector", "#selector"})
		})

		ItShouldMakeAnElementRequest("elements", "POST", `{"using": "css selector", "value": "#selector"}`)

		It("should return a slice of elements with IDs and sessions", func() {
			Expect(elements[0].ID).To(Equal("some-id"))
			Expect(elements[0].Session).To(Equal(session))
			Expect(elements[1].ID).To(Equal("some-other-id"))
			Expect(elements[1].Session).To(Equal(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetElements(Selector{"css selector", "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetText", func() {
		var text string

		BeforeEach(func() {
			bus.SendCall.Result = `"some text"`
			text, err = element.GetText()
		})

		ItShouldMakeAnElementRequest("text", "GET")

		It("should return the visible text on the element", func() {
			Expect(text).To(Equal("some text"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error indicating the bus failed to retrieve the text", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetText()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetName", func() {
		var text string

		BeforeEach(func() {
			bus.SendCall.Result = `"some-name"`
			text, err = element.GetName()
		})

		ItShouldMakeAnElementRequest("name", "GET")

		It("should return the tag name of the element", func() {
			Expect(text).To(Equal("some-name"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error indicating the bus failed to retrieve the tag name", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetName()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetAttribute", func() {
		var value string

		BeforeEach(func() {
			bus.SendCall.Result = `"some value"`
			value, err = element.GetAttribute("some-name")
		})

		ItShouldMakeAnElementRequest("attribute/some-name", "GET")

		It("should return the value of the attribute", func() {
			Expect(value).To(Equal("some value"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetAttribute("some-name")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetCSS", func() {
		var value string

		BeforeEach(func() {
			bus.SendCall.Result = `"some value"`
			value, err = element.GetCSS("some-property")
		})

		ItShouldMakeAnElementRequest("css/some-property", "GET")

		It("should return the value of the CSS property", func() {
			Expect(value).To(Equal("some value"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetCSS("some-property")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Click", func() {
		BeforeEach(func() {
			err = element.Click()
		})

		ItShouldMakeAnElementRequest("click", "POST")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(element.Click()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Clear", func() {
		BeforeEach(func() {
			err = element.Clear()
		})

		ItShouldMakeAnElementRequest("clear", "POST")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(element.Clear()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Value", func() {
		BeforeEach(func() {
			err = element.Value("text")
		})

		ItShouldMakeAnElementRequest("value", "POST", `{"value": ["t", "e", "x", "t"]}`)

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(element.Value("text")).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsSelected", func() {
		var value bool

		BeforeEach(func() {
			bus.SendCall.Result = "true"
			value, err = element.IsSelected()
		})

		ItShouldMakeAnElementRequest("selected", "GET")

		It("should return the selected status", func() {
			Expect(value).To(BeTrue())
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.IsSelected()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsDisplayed", func() {
		var value bool

		BeforeEach(func() {
			bus.SendCall.Result = "true"
			value, err = element.IsDisplayed()
		})

		ItShouldMakeAnElementRequest("displayed", "GET")

		It("should return the displayed status", func() {
			Expect(value).To(BeTrue())
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.IsDisplayed()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsEnabled", func() {
		var value bool

		BeforeEach(func() {
			bus.SendCall.Result = "true"
			value, err = element.IsEnabled()
		})

		ItShouldMakeAnElementRequest("enabled", "GET")

		It("should return the enabled status", func() {
			Expect(value).To(BeTrue())
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.IsEnabled()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Submit", func() {
		BeforeEach(func() {
			err = element.Submit()
		})

		ItShouldMakeAnElementRequest("submit", "POST")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(element.Submit()).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsEqualTo", func() {
		var (
			equal        bool
			otherElement *Element
		)

		BeforeEach(func() {
			otherElement = &Element{"other-id", session}
			bus.SendCall.Result = "true"
			equal, err = element.IsEqualTo(otherElement)
		})

		ItShouldMakeAnElementRequest("equals/other-id", "GET")

		It("should return whether the elements are equal", func() {
			Expect(equal).To(BeTrue())
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.IsEqualTo(otherElement)
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
