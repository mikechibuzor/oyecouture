-- Lets a fulfillment event record the out-for-delivery stage.
--
-- Migration 000018 added this stage and widened orders_fulfillment_stage_valid,
-- but not the matching constraint on order_fulfillment_events, which was
-- written in 000016 and lists the stages separately. A transition writes both:
-- it updates the order and inserts the event in one transaction. The event
-- insert failed the CHECK, the whole transaction rolled back, and the admin saw
-- only "order tracking could not be loaded" - so moving an order to out for
-- delivery has been impossible in production since the stage shipped.
--
-- The same shape of bug as the email type in 000020: one stage, three separate
-- CHECK constraints, and only some of them widened. constraints_test.go now
-- asserts every stage against both of the stage constraints, which is what
-- stops the next one reaching production.
ALTER TABLE order_fulfillment_events
    DROP CONSTRAINT IF EXISTS order_fulfillment_events_stage_valid;

ALTER TABLE order_fulfillment_events
    ADD CONSTRAINT order_fulfillment_events_stage_valid CHECK (
        stage IN (
            'payment_confirmed',
            'measurements_approved',
            'preparing_order',
            'in_production',
            'quality_check',
            'ready_for_dispatch',
            'dispatched',
            'out_for_delivery',
            'delivered'
        )
    );
