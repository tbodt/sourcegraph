import React from 'react'
import { CompletionItem } from 'sourcegraph'
import { CompletionWidget } from './CompletionWidget'
import { mount } from 'enzyme'

const COMPLETION_ITEM_2: CompletionItem = { label: 'b' }

describe('CompletionWidgetDropdown', () => {
    test('simple', () => {
        const textArea = document.createElement('textarea')
        expect(
            mount(
                <CompletionWidget
                    completionListOrError={{ items: [{ label: 'a' }, COMPLETION_ITEM_2] }}
                    onSelectItem={() => undefined}
                    textArea={textArea}
                    listClassName="list-class-name"
                    listItemClassName="list-item-class-name"
                    loadingClassName="loading-class-name"
                    noResultsClassName="no-results-class-name"
                    selectedListItemClassName="selected-list-item-class-name"
                    widgetClassName="widget-class-name"
                    widgetContainerClassName="widget-container-class-name"
                />
            ).children()
        ).toMatchSnapshot()
    })
})
